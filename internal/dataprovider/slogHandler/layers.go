package slogHandler

import (
	"log/slog"
	"slices"
)

type handlerLayer struct {
	group string
	attrs []slog.Attr
}

func layersWithAttrs(layers []handlerLayer, attrs []slog.Attr) []handlerLayer {
	if len(attrs) == 0 {
		return layers
	}

	if len(layers) == 0 {
		return []handlerLayer{
			{
				attrs: attrs,
			},
		}
	}

	lCopy := make([]handlerLayer, len(layers))
	copy(lCopy, layers)

	lastLayerIndex := len(lCopy) - 1
	lCopy[lastLayerIndex].attrs = append(slices.Clip(lCopy[lastLayerIndex].attrs), attrs...)

	return lCopy
}

func layersWithGroup(layers []handlerLayer, group string) []handlerLayer {
	if group == "" {
		return layers
	}

	if len(layers) == 0 {
		return []handlerLayer{
			{
				group: group,
			},
		}
	}

	lCopy := make([]handlerLayer, len(layers)+1)
	copy(lCopy, layers)

	lastLayerIndex := len(lCopy) - 1
	lCopy[lastLayerIndex] = handlerLayer{
		group: group,
	}

	return lCopy
}

func layersToAttrs(layers []handlerLayer) []slog.Attr {
	if len(layers) == 0 {
		return nil
	}

	layer := layers[0]
	attrs := slices.Clip(layer.attrs)

	if len(layers) > 1 {
		attrs = append(attrs, layersToAttrs(layers[1:])...)
	}

	if layer.group != "" {
		return []slog.Attr{
			{
				Key:   layer.group,
				Value: slog.GroupValue(attrs...),
			},
		}
	}

	return attrs
}
