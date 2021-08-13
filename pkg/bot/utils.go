package bot

import "errors"

func ParseChannelMode(num int) (err error, mode ChannelMode) {
	mode = ChannelMode(num)
	if mode < ChannelModeNormal || mode >= ChannelModeBoundary {
		return errors.New("mode value out of bounds"), mode
	}
	return nil, mode
}
