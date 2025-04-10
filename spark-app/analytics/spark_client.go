package analytics

import (
	"context"
	"fmt"
	"log"
	"time"

	"spark-app/config"

	"github.com/apache/spark-connect-go/v35/spark/sql"
	"github.com/apache/spark-connect-go/v35/spark/sql/utils"
)

type Recommendation struct {
	Hour  int64 `json:"hour"`
	Count int64 `json:"count"`
}

func AnalyzeCreatedAtHours() ([]Recommendation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	// Create a new Spark session
	spark, err := sql.NewSessionBuilder().
		Remote(config.SparkURL()).
		Build(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to create Spark session: %w", err)
	}
	defer utils.WarnOnError(spark.Stop, func(err error) {
		log.Printf("Warning: error while stopping Spark session: %v", err)
	})

	_, err = spark.Sql(ctx, "SET spark.hadoop.fs.defaultFS=hdfs://hadoop-namenode:9000")
	if err != nil {
		return nil, fmt.Errorf("failed to set HDFS configuration: %w", err)
	}

	// Correct HDFS path format
	hdfsPath := fmt.Sprintf("hdfs://hadoop-namenode:9000%s/*.jsonl", config.HDFSDataPath())

	// Form the SQL query
	query := fmt.Sprintf(`
		SELECT
			HOUR(CAST(created_at AS TIMESTAMP)) AS hour,
			COUNT(*) AS count
		FROM
			JSON.`+"`%s`"+`
		GROUP BY
			HOUR(CAST(created_at AS TIMESTAMP))
		ORDER BY
			count DESC
	`, hdfsPath)

	// Execute the query
	df, err := spark.Sql(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to execute SQL query: %w", err)
	}

	// Collect results
	rows, err := df.Collect(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to collect results: %w", err)
	}

	// Process results
	var recommendations []Recommendation
	for _, row := range rows {
		values := row.Values()
		if len(values) < 2 {
			log.Printf("Warning: row has fewer values than expected: %v", values)
			continue
		}

		hour, ok := values[0].(int64)
		if !ok {
			log.Printf("Warning: unexpected hour type: %T", values[0])
			continue
		}

		count, ok := values[1].(int64)
		if !ok {
			log.Printf("Warning: unexpected count type: %T", values[1])
			continue
		}

		recommendations = append(recommendations, Recommendation{
			Hour:  hour,
			Count: count,
		})
	}

	return recommendations, nil
}
