package ws

/**
 * This file is part of Gourd.
 *
 * @link     http://gourd.kyour.cn
 * @document http://gourd.kyour.cn/doc
 * @contact  kyour@vip.qq.com
 * @license  https://https://github.com/kyour-cn/gourd/blob/master/LICENSE
 */

//ws消息类型常量
//内容来自 github.com/gorilla/websocket/Conn.go
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	//text message表示文本数据消息。文本消息负载是
	//解释为UTF-8编码文本数据。
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	//BinaryMessage表示二进制数据消息。
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	//CloseMessage表示关闭控制消息。可选消息
	//有效负载包含数字代码和文本。使用FormatCloseMessage
	//函数设置关闭消息有效负载的格式。
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	//PingMessage表示ping控制消息。可选消息负载
	//是UTF-8编码文本。
	PingMessage = 9

	// PongMessage denotes a pong control message. The optional message payload
	// is UTF-8 encoded text.
	//PongMessage表示pong控制消息。可选消息负载
	//是UTF-8编码文本。
	PongMessage = 10
)
