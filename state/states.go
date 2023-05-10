/*
 * Copyright (c) Zinnatullin Chingiz, 2022.
 */

package state

// FuncState хранит текущее состояние работы с пользователем
var FuncState = intInt{
	m: make(map[int]int),
}
