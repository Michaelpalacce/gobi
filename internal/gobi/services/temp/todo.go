package temp

// func (p *Processor) processItemSaveMessage(websocketMessage messages.WebsocketMessage) error {
// 	var (
// 		itemSavePayload v1.ItemSavePayload
// 		err             error
// 	)
//
// 	if err = json.Unmarshal(websocketMessage.Payload, &itemSavePayload); err != nil {
// 		return err
// 	}
// 	item := itemSavePayload.Item
//
// 	if item.SHA256 == p.WebsocketClient.StorageDriver.CalculateSHA256(item) {
// 		// @TODO: touch the files in peer servers, send an event via redis
// 		if err := p.WebsocketClient.StorageDriver.Touch(item); err != nil {
// 			return err
// 		}
// 		slog.Info("Item already exists locally", "item", item)
// 		return nil
// 	}
//
// 	syncStrategy := p.SyncStrategy
// 	if err := syncStrategy.FetchSingle(item, storage.ConflictModeYes); err != nil {
// 		return err
// 	}
//
// 	return nil
// }

// // SendItem will send an item to the the client/server
// func (c *WebsocketClient) SendItem(item models.Item) error {
// 	slog.Debug("Sending file to server", "item", item)
//
// 	reader, err := c.StorageDriver.GetReader(item)
// 	if err != nil {
// 		return err
// 	}
// 	defer reader.Close()
//
// 	buffer := make([]byte, 1024)
//
// 	for {
// 		n, err := reader.Read(buffer)
// 		if err == io.EOF {
// 			break
// 		}
//
// 		if err != nil {
// 			return fmt.Errorf("error reading: %w", err)
// 		}
//
// 		err = c.Conn.WriteMessage(websocket.BinaryMessage, buffer[:n])
// 		if err != nil {
// 			return fmt.Errorf("error reading file chunk: %w", err)
// 		}
// 	}
//
// 	slog.Info("File Sent Successfully", "item", item.ServerPath, "vault", item.VaultName)
//
// 	return nil
// }
//
// // FetchItem will receive an item from the client/server
// // @NOTE: You must tell the client/server to start sending the file first
// func (c *WebsocketClient) FetchItem(item models.Item) error {
// 	c.Conn.SetReadDeadline(time.Now().Add(30 * time.Second))
// 	defer c.Conn.SetReadDeadline(time.Time{})
//
// 	slog.Debug("Fetching file", "item", item)
//
// 	writer, err := c.StorageDriver.GetWriter(item)
// 	if err != nil {
// 		return err
// 	}
//
// 	defer func() {
// 		writer.Close()
// 	}()
//
// 	if item.Size == 0 {
// 		slog.Info("File Fetched Successfully", "item", item)
// 		return nil
// 	}
//
// 	bytesRead := 0
// 	for {
// 		messageType, message, err := c.Conn.ReadMessage()
// 		if err != nil {
// 			return err
// 		}
//
// 		if messageType != websocket.BinaryMessage {
// 			return fmt.Errorf("invalid messageType received: %d, expected 2 (BinaryMessage)", messageType)
// 		}
//
// 		writer.Write(message)
//
// 		bytesRead += len(message)
// 		if bytesRead == item.Size {
// 			writer.Close()
// 			break
// 		}
//
// 		if bytesRead > item.Size {
// 			return fmt.Errorf("expected %d bytes, but got %d", item.Size, bytesRead)
// 		}
// 	}
//
// 	slog.Info("File Fetched Successfully", "Item", item.ServerPath, "vault", item.VaultName)
//
// 	return nil
// }
