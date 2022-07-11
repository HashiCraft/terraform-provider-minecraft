package minecraft

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/seeruk/minecraft-rcon/rcon"
)

type Client struct {
	client *rcon.Client
}

func New(address string, password string) (*Client, error) {
	addressParts := strings.Split(address, ":")
	host := addressParts[0]
	port, err := strconv.Atoi(addressParts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid port %s", addressParts[1])
	}

	client, err := rcon.NewClient(host, port, password)
	if err != nil {
		return nil, err
	}

	return &Client{client}, nil
}

// Creates a block.
func (c Client) CreateBlock(ctx context.Context, material string, x, y, z int) error {
	command := fmt.Sprintf("setblock %d %d %d %s replace", x, y, z, material)
	_, err := c.client.SendCommand(command)
	if err != nil {
		return err
	}

	return nil
}

// Deletes a block.
func (c Client) DeleteBlock(ctx context.Context, x, y, z int) error {
	command := fmt.Sprintf("setblock %d %d %d minecraft:air replace", x, y, z)
	_, err := c.client.SendCommand(command)
	if err != nil {
		return err
	}

	return nil
}

// Creates an entity.
func (c Client) CreateEntity(ctx context.Context, entity string, position string, id string) error {
	command := fmt.Sprintf("summon minecraft:%s %s {CustomName:'{\"text\":\"%s\"}'}", entity, position, id)
	_, err := c.client.SendCommand(command)
	if err != nil {
		return err
	}

	return nil
}

// Deletes an entity.
func (c Client) DeleteEntity(ctx context.Context, entity string, position string, id string) error {
	// Remove the entity.
	command := fmt.Sprintf("kill @e[type=minecraft:%s,nbt={CustomName:'{\"text\":\"%s\"}'}]", entity, id)
	_, err := c.client.SendCommand(command)
	if err != nil {
		return err
	}

	// Remove the entity from inventories.
	command = fmt.Sprintf("clear @a minecraft:%s{display:{Name:'{\"text\":\"%s\"}'}}", entity, id)
	_, err = c.client.SendCommand(command)
	if err != nil {
		return err
	}

	return nil
}
