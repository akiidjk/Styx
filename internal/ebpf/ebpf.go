package ebpf

import (
	"fmt"

	ebpfGenerated "github.com/akiidjk/styx/internal/ebpf/generated"
	l "github.com/akiidjk/styx/internal/utils/logger"
	"github.com/cilium/ebpf"
)

var logger = l.GetLogger()

var objects = []string{"FilterIp", "FilterMac", "Dispacher", "FilterPort"}
var ObjectsEBPFMap = make(map[string]interface{})

/*
Load all the eBPF objects
*/
func LoadObjectsEBPF() {
	logger.Info().Msg("Loading eBPF objects")
	for index, object := range objects {
		logger.Debug().Str("object", object).Msg("Loading eBPF object")
		ObjectsEBPFMap[object] = make(map[string]interface{})
		switch object {
		case "FilterIp":
			ObjectsEBPFMap[object] = LoadEBPFObjectsFilterIp()
		case "FilterMac":
			ObjectsEBPFMap[object] = loadEBPFObjectsFilterMac()
		case "Dispacher":
			ObjectsEBPFMap[object] = LoadEBPFObjectsDispacher()
		default:
			logger.Warn().Str("object", object).Msg("Unknown eBPF object type")
			continue
		}
		logger.Info().
			Str("object", object).
			Int("mappings", index).
			Msg("Successfully loaded eBPF object")
	}

	if len(ObjectsEBPFMap) != len(objects) {
		logger.Fatal().
			Int("loaded_objects", len(ObjectsEBPFMap)).
			Int("total_objects", len(objects)).
			Msg("Failed to load all eBPF objects")
		return
	}

	if len(ObjectsEBPFMap) == 0 {
		logger.Warn().Msg("No eBPF objects loaded")
	}

	logger.Info().
		Int("total_objects", len(ObjectsEBPFMap)).
		Msg("Completed loading all eBPF objects")
}

/*
Get all the programs from the eBPF objects
*/
func GetPrograms() []*ebpf.Program {
	logger.Debug().Msg("Getting programs from ObjectsEBPFMap")
	var loadedPrograms []*ebpf.Program

	if filterIpObj, ok := ObjectsEBPFMap["FilterIp"].(*ebpfGenerated.FilteripObjects); ok {
		loadedPrograms = append(loadedPrograms, filterIpObj.XdpFilterIp)
	} else {
		logger.Fatal().Msgf("Failed to assert FilterIp object. Value: %v", ObjectsEBPFMap["FilterIp"])
		return nil
	}

	if filterMacObj, ok := ObjectsEBPFMap["FilterMac"].(*ebpfGenerated.FiltermacObjects); ok {
		loadedPrograms = append(loadedPrograms, filterMacObj.XdpFilterMac)
	} else {
		logger.Fatal().Msgf("Failed to assert FilterMac object. Value: %v", ObjectsEBPFMap["FilterMac"])
		return nil
	}

	return loadedPrograms
}

/*
Get a specific object from the ObjectsEBPFMap

Use example:

	obj, err := GetObject[*ebpfGenerated.DispacherObjects]("Dispacher")
	if err != nil {
		logger.Fatal().Err(err).Msg("Failed to get Dispacher object")
	}
	addPrograms(GetPrograms(), (*obj).ProgArray)
*/
func GetObject[T any](key string) (*T, error) {
	logger.Debug().Str("key", key).Msg("Getting object from ObjectsEBPFMap")
	obj, exists := ObjectsEBPFMap[key]
	if !exists || obj == nil {
		logger.Fatal().Msg(fmt.Sprintf("failed to find %s in ObjectsEBPFMap", key))
		return nil, fmt.Errorf("Failed to find %s in ObjectsEBPFMap", key)
	}

	switch typedObj := obj.(type) {
	case *T:
		return typedObj, nil
	case T:
		return &typedObj, nil
	default:
		logger.Fatal().Msg(fmt.Sprintf("Failed to assert %s as type %T", key, new(T)))
		return nil, fmt.Errorf("Failed to assert %s as type %T", key, new(T))
	}
}

func CloseAllObject() {
	logger.Info().Msg("Starting to close all eBPF objects")

	for name, object := range ObjectsEBPFMap {
		if object == nil {
			logger.Warn().
				Str("object", name).
				Msg("Object is nil, skipping")
			continue
		}

		logger.Debug().
			Str("object", name).
			Msg("Attempting to close object")

		switch name {
		case "FilterIp":
			obj, ok := object.(ebpfGenerated.FilteripObjects)
			if !ok {
				logger.Error().
					Str("object", name).
					Msg("Failed to cast object to FilteripObjects")
				continue
			}
			if err := obj.Close(); err != nil {
				logger.Error().
					Str("object", name).
					Err(err).
					Msg("Error closing FilterIp object")
			} else {
				logger.Info().
					Str("object", name).
					Msg("Successfully closed FilterIp object")
			}

		case "FilterMac":
			obj, ok := object.(ebpfGenerated.FiltermacObjects)
			if !ok {
				logger.Error().
					Str("object", name).
					Msg("Failed to cast object to FiltermacObjects")
				continue
			}
			if err := obj.Close(); err != nil {
				logger.Error().
					Str("object", name).
					Err(err).
					Msg("Error closing FilterMac object")
			} else {
				logger.Info().
					Str("object", name).
					Msg("Successfully closed FilterMac object")
			}

		default:
			logger.Warn().
				Str("object", name).
				Msg("Unknown object type")
		}
	}

	logger.Info().Msg("Finished closing all eBPF objects")
}
