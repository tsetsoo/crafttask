package crafttask

import "errors"

var errBlockDoesNotExist = errors.New("block does not exist")
var errParentBlockDoesNotExist = errors.New("parent block does not exist")
var errBlockMovedToItsChild = errors.New("block attempt to move to its child")
