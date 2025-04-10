package cli

import (
	"flag"
	"fmt"

	"banned-products-system/internal/storage"
)

// Handler обрабатывает CLI-команды
type Handler struct {
	store *storage.Storage
}

// NewHandler создает новый обработчик CLI-команд
func NewHandler(store *storage.Storage) *Handler {
	return &Handler{
		store: store,
	}
}

// HandleCommand обрабатывает команду CLI
func (h *Handler) HandleCommand(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("команда не указана")
	}

	cmd := args[0]
	remainingArgs := args[1:]

	switch cmd {
	case "add":
		return h.handleAddCommand(remainingArgs)
	case "remove":
		return h.handleRemoveCommand(remainingArgs)
	case "update":
		return h.handleUpdateCommand(remainingArgs)
	case "list":
		return h.handleListCommand()
	default:
		return fmt.Errorf("неизвестная команда: %s", cmd)
	}
}

// handleAddCommand обрабатывает команду добавления товара
func (h *Handler) handleAddCommand(args []string) error {
	addCmd := flag.NewFlagSet("add", flag.ExitOnError)
	sku := addCmd.String("sku", "", "SKU товара для добавления в список запрещенных")
	reason := addCmd.String("reason", "", "Причина запрета товара")

	if err := addCmd.Parse(args); err != nil {
		return err
	}

	if *sku == "" {
		return fmt.Errorf("требуется указать SKU товара с помощью флага -sku")
	}

	if *reason == "" {
		return fmt.Errorf("требуется указать причину запрета с помощью флага -reason")
	}

	if err := h.store.Add(*sku, *reason); err != nil {
		return err
	}

	fmt.Printf("Товар с SKU %s добавлен в список запрещенных\n", *sku)
	return nil
}

// handleRemoveCommand обрабатывает команду удаления товара
func (h *Handler) handleRemoveCommand(args []string) error {
	removeCmd := flag.NewFlagSet("remove", flag.ExitOnError)
	sku := removeCmd.String("sku", "", "SKU товара для удаления из списка запрещенных")

	if err := removeCmd.Parse(args); err != nil {
		return err
	}

	if *sku == "" {
		return fmt.Errorf("требуется указать SKU товара с помощью флага -sku")
	}

	if err := h.store.Remove(*sku); err != nil {
		return err
	}

	fmt.Printf("Товар с SKU %s удален из списка запрещенных\n", *sku)
	return nil
}

// handleUpdateCommand обрабатывает команду обновления товара
func (h *Handler) handleUpdateCommand(args []string) error {
	updateCmd := flag.NewFlagSet("update", flag.ExitOnError)
	sku := updateCmd.String("sku", "", "SKU товара для обновления в списке запрещенных")
	reason := updateCmd.String("reason", "", "Новая причина запрета товара")

	if err := updateCmd.Parse(args); err != nil {
		return err
	}

	if *sku == "" {
		return fmt.Errorf("требуется указать SKU товара с помощью флага -sku")
	}

	if *reason == "" {
		return fmt.Errorf("требуется указать новую причину запрета с помощью флага -reason")
	}

	if err := h.store.Update(*sku, *reason); err != nil {
		return err
	}

	fmt.Printf("Причина запрета товара с SKU %s обновлена\n", *sku)
	return nil
}

// handleListCommand обрабатывает команду вывода списка запрещенных товаров
func (h *Handler) handleListCommand() error {
	products := h.store.List()

	fmt.Println("Список запрещенных товаров:")
	fmt.Println("===========================")

	if len(products) == 0 {
		fmt.Println("Список пуст")
		return nil
	}

	for i, product := range products {
		fmt.Printf("%d. SKU: %s\n   Причина: %s\n   Дата добавления: %s\n\n",
			i+1, product.SKU, product.Reason, product.BannedAt.Format("2006-01-02 15:04:05"))
	}

	return nil
}
