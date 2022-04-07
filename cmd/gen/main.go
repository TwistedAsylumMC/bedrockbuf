package main

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/sandertv/gophertunnel/minecraft/protocol"
	"github.com/sandertv/gophertunnel/minecraft/protocol/packet"
	"image/color"
	"os"
	"reflect"
	"regexp"
	"sort"
	"strings"
	"text/template"
)

var (
	mapRegex      = regexp.MustCompile(`map\[(.+)](\S+)`)
	protocolTypes = []any{
		&RepeatedRGBA{},
		&TransferStackRequestAction{},
		&Vector2f{},
		&Vector2i{},
		&Vector3f{},
		&Vector3i{},
		&color.RGBA{},
		&protocol.Attribute{},
		&protocol.BehaviourPackInfo{},
		&protocol.BlockChangeEntry{},
		&protocol.BlockEntry{},
		&protocol.CacheBlob{},
		&protocol.Command{},
		&protocol.CommandEnum{},
		&protocol.CommandEnumConstraint{},
		&protocol.CommandOrigin{},
		&protocol.CommandOutputMessage{},
		&protocol.CommandOverload{},
		&protocol.CommandParameter{},
		&protocol.CreativeItem{},
		&protocol.EducationExternalLinkSettings{},
		&protocol.EducationSharedResourceURI{},
		&protocol.EnchantmentInstance{},
		&protocol.EnchantmentOption{},
		&protocol.EntityLink{},
		&protocol.ExperimentData{},
		&protocol.GameRule{},
		&protocol.InventoryAction{},
		&protocol.ItemComponentEntry{},
		&protocol.ItemEnchantments{},
		&protocol.ItemEntry{},
		&protocol.ItemInstance{},
		&protocol.ItemStack{},
		&protocol.ItemStackRequest{},
		&protocol.ItemStackResponse{},
		&protocol.ItemType{},
		&protocol.LegacySetItemSlot{},
		&protocol.MapDecoration{},
		&protocol.MapTrackedObject{},
		&protocol.MaterialReducer{},
		&protocol.MaterialReducerOutput{},
		&protocol.PersonaPiece{},
		&protocol.PersonaPieceTintColour{},
		&protocol.PlayerBlockAction{},
		&protocol.PlayerListEntry{},
		&protocol.PlayerMovementSettings{},
		&protocol.PotionContainerChangeRecipe{},
		&protocol.PotionRecipe{},
		&protocol.RecipeIngredientItem{},
		&protocol.ScoreboardEntry{},
		&protocol.ScoreboardIdentityEntry{},
		&protocol.Skin{},
		&protocol.SkinAnimation{},
		&protocol.StackRequestSlotInfo{},
		&protocol.StackResourcePack{},
		&protocol.StackResponseContainerInfo{},
		&protocol.StackResponseSlotInfo{},
		&protocol.StructureSettings{},
		&protocol.SubChunkEntry{},
		&protocol.TexturePackInfo{},
	}
	interfaceTypes = map[string][]any{
		"EventData": {
			&protocol.AchievementAwardedEventData{},
			&protocol.AgentCommandEventData{},
			&protocol.BellUsedEventData{},
			&protocol.BossKilledEventData{},
			&protocol.CauldronInteractEventData{},
			&protocol.CauldronUsedEventData{},
			&protocol.ComposterInteractEventData{},
			&protocol.EntityDefinitionTriggerEventData{},
			&protocol.EntityInteractEventData{},
			&protocol.ExtractHoneyEventData{},
			&protocol.FishBucketedEventData{},
			&protocol.MobBornEventData{},
			&protocol.MobKilledEventData{},
			&protocol.MovementAnomalyEventData{},
			&protocol.MovementCorrectedEventData{},
			&protocol.PatternRemovedEventData{},
			&protocol.PetDiedEventData{},
			&protocol.PlayerDiedEventData{},
			&protocol.PortalBuiltEventData{},
			&protocol.PortalUsedEventData{},
			&protocol.RaidUpdateEventData{},
			&protocol.SlashCommandExecutedEventData{},
		},
		"InventoryTransactionData": {
			&protocol.MismatchTransactionData{},
			&protocol.NormalTransactionData{},
			&protocol.ReleaseItemTransactionData{},
			&protocol.UseItemOnEntityTransactionData{},
			&protocol.UseItemTransactionData{},
		},
		"Recipe": {
			&protocol.FurnaceDataRecipe{},
			&protocol.FurnaceRecipe{},
			&protocol.MultiRecipe{},
			&protocol.ShapedChemistryRecipe{},
			&protocol.ShapedRecipe{},
			&protocol.ShapelessChemistryRecipe{},
			&protocol.ShapelessRecipe{},
			&protocol.ShulkerBoxRecipe{},
		},
		"StackRequestAction": {
			&protocol.AutoCraftRecipeStackRequestAction{},
			&protocol.BeaconPaymentStackRequestAction{},
			&protocol.ConsumeStackRequestAction{},
			&protocol.CraftCreativeStackRequestAction{},
			&protocol.CraftGrindstoneRecipeStackRequestAction{},
			&protocol.CraftLoomRecipeStackRequestAction{},
			&protocol.CraftNonImplementedStackRequestAction{},
			&protocol.CraftRecipeOptionalStackRequestAction{},
			&protocol.CraftRecipeStackRequestAction{},
			&protocol.CraftResultsDeprecatedStackRequestAction{},
			&protocol.CreateStackRequestAction{},
			&protocol.DestroyStackRequestAction{},
			&protocol.DropStackRequestAction{},
			&protocol.LabTableCombineStackRequestAction{},
			&protocol.MineBlockStackRequestAction{},
			&protocol.PlaceInContainerStackRequestAction{},
			&protocol.PlaceStackRequestAction{},
			&protocol.SwapStackRequestAction{},
			&protocol.TakeOutContainerStackRequestAction{},
			&protocol.TakeStackRequestAction{},
		},
	}
)

type Messages struct {
	Messages []Message
}

type Message struct {
	OneOf  string
	Name   string
	Fields []Field
}

type Field struct {
	Prepend string
	Type    string
	Name    string
	Index   int
}

type RepeatedRGBA struct {
	RGBA []color.RGBA
}

type TransferStackRequestAction struct {
	Count               byte
	Source, Destination protocol.StackRequestSlotInfo
}

type Vector2f struct {
	X, Y float32
}

type Vector2i struct {
	X, Y int32
}

type Vector3f struct {
	X, Y, Z float32
}

type Vector3i struct {
	X, Y, Z int32
}

func main() {
	pool := packet.NewPool()
	var packetMessages []Message
	var typeMessages []Message
	for _, pkFunc := range pool {
		message, typeMessage := createMessages(pkFunc())
		packetMessages = append(packetMessages, message)
		typeMessages = append(typeMessages, typeMessage...)
	}
	err := executeTemplate("packets", packetMessages)
	if err != nil {
		panic(err)
	}

	for _, x := range protocolTypes {
		message, typeMessage := createMessages(x)
		typeMessages = append(typeMessages, message)
		typeMessages = append(typeMessages, typeMessage...)
	}
	for _, types := range interfaceTypes {
		for _, x := range types {
			message, typeMessage := createMessages(x)
			typeMessages = append(typeMessages, message)
			typeMessages = append(typeMessages, typeMessage...)
		}
	}
	err = executeTemplate("types", typeMessages)
	if err != nil {
		panic(err)
	}
}

func executeTemplate(name string, messages []Message) error {
	tmpl := template.Must(template.ParseFiles(fmt.Sprintf("cmd/gen/%s.proto.tmpl", name)))
	out, err := os.Create(fmt.Sprintf("%s.proto", name))
	if err != nil {
		panic(err)
	}
	sort.Slice(messages, func(i, j int) bool {
		return messages[i].Name < messages[j].Name
	})
	err = tmpl.Execute(out, Messages{messages})
	if err != nil {
		panic(err)
	}
	return nil
}

func createMessages(x interface{}) (message Message, typesMessages []Message) {
	name := fmt.Sprintf("%T", x)
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[1]
	}
	if _, ok := x.(packet.Packet); ok {
		name += "Packet"
	}
	message.Name = name
	val := reflect.ValueOf(x).Elem()
	interfaceIndex := -1
	for i := 0; i < val.NumField(); i++ {
		fieldVal := val.Type().Field(i)
		fieldName := strcase.ToSnake(fieldVal.Name)
		fieldTypeName, prepend := convertType(fieldVal.Type.String(), false)
		if _, ok := interfaceTypes[fieldTypeName]; ok {
			interfaceIndex = i
			fieldTypeName = "I" + fieldTypeName
		}
		field := Field{
			Type:  fieldTypeName,
			Name:  fieldName,
			Index: i + 1,
		}
		if len(prepend) > 0 {
			field.Prepend = strings.Join(prepend, " ") + " "
		}
		message.Fields = append(message.Fields, field)
	}
	if interfaceIndex >= 0 {
		field := message.Fields[interfaceIndex]
		typeMessage := Message{OneOf: field.Name, Name: field.Type}
		types := interfaceTypes[field.Type[1:]]
		for i, x := range types {
			fieldType := reflect.ValueOf(x).Elem().Type().String()
			fieldName := strcase.ToSnake(fieldType)
			fieldTypeName, _ := convertType(fieldType, false)
			typeMessage.Fields = append(typeMessage.Fields, Field{
				Type:  fieldTypeName,
				Name:  fieldName,
				Index: i + 1,
			})
		}
		typesMessages = append(typesMessages, typeMessage)
	}
	return
}

func convertType(x string, fromSlice bool) (string, []string) {
	if x == "[]byte" {
		return "bytes", nil
	} else if fromSlice && x == "[]color.RGBA" {
		return "RepeatedRGBA", nil
	} else if strings.HasPrefix(x, "[]") {
		sliceType, prepend := convertType(x[2:], true)
		return sliceType, append([]string{"repeated"}, prepend...)
	} else if strings.HasPrefix(x, "map") {
		groups := mapRegex.FindStringSubmatch(x)
		key, _ := convertType(groups[1], false)
		value, prepend := convertType(groups[2], false)
		var prependStr string
		if len(prepend) > 0 {
			prependStr = strings.Join(prepend, " ") + " "
		}
		return fmt.Sprintf("map<%s, %s%s>", key, prependStr, value), nil
	} else if strings.HasPrefix(x, "*") {
		ptrType, prepend := convertType(x[1:], false)
		return ptrType, append([]string{"optional"}, prepend...)
	} else if strings.Contains(x, ".") {
		return convertType(strings.Split(x, ".")[1], false)
	}
	switch x {
	case "[3]int8":
		return "Vector3i", nil
	case "byte", "uint8", "uint16":
		return "uint32", nil
	case "int8", "int16":
		return "int32", nil
	case "float32":
		return "float", nil
	case "float64":
		return "double", nil
	case "interface", "interface {}":
		return "google.protobuf.Any", nil

	case "BlockPos", "SubChunkPos":
		return "Vector3i", nil
	case "ChunkPos":
		return "Vector2i", nil
	case "transferStackRequestAction":
		return "TransferStackRequestAction", nil
	case "UUID":
		return "string", nil
	case "Vec2":
		return "Vector2f", nil
	case "Vec3":
		return "Vector3f", nil
	default:
		return x, nil
	}
}
