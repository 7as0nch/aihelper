/* *
 * @Author: chengjiang
 * @Date: 2025-10-04 17:16:55
 * @Description:
**/

package models

type Message struct {
	Model
	RoomID  int64  `json:"room_id"`
	UserID  int64  `json:"user_id"`
	Content string `json:"content"`
}
