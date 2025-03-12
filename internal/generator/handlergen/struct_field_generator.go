package handlergen

import (
	"strings"

	"github.com/mrbryside/go-generate/internal/utils/myfile"
	"github.com/mrbryside/go-generate/internal/utils/mymap"
	"github.com/mrbryside/go-generate/internal/utils/mystr"
)

// Helper function to generate struct fields from request map
func generateStructFields(
	payloadType string,
	handlerName string,
	requestOrResponse *mymap.OrderedMap,
	fields []string,
	nss []myfile.NewStruct,
	previousTypeName string,
) (string, []myfile.NewStruct) {
	iter := requestOrResponse.EntriesIter()
	for {
		pair, ok := iter()
		if !ok {
			break
		}
		fieldName := pair.Key
		fieldValue := pair.Value
		var baseType string
		var required bool

		switch v := fieldValue.(type) {
		case string:
			fieldReplace := strings.Replace(v, "|", " ", 1)
			typeParts := strings.Fields(fieldReplace)
			baseType = typeParts[0]
			required = len(typeParts) > 1 && typeParts[1] == "required"

			camelCaseFieldName := mystr.CamelCase(fieldName)
			exportedFieldName := mystr.CapitalizeFirstLetter(camelCaseFieldName)

			tags := "`json:\"" + fieldName + "\""
			if required {
				tags += " validate:\"required\""
			}
			tags += "`"

			fields = append(fields, exportedFieldName+" "+baseType+" "+tags)
		case *mymap.OrderedMap:
			fieldReplace := strings.Replace(fieldName, "|", " ", 1)
			fieldParts := strings.Fields(fieldReplace)
			fieldName := fieldParts[0]
			required = len(fieldParts) > 1 && fieldParts[1] == "required"
			camelCaseFieldName := mystr.CamelCase(fieldName)
			exportedFieldName := mystr.CapitalizeFirstLetter(camelCaseFieldName)
			tags := "`json:\"" + fieldName + "\""
			if required {
				tags += " validate:\"required\""
			}
			tags += "`"

			// using previousTypeName to generate recursive struct name
			typeName := ""
			if previousTypeName == "" {
				typeName = handlerName + exportedFieldName + payloadType
			} else {
				typeName = strings.Replace(previousTypeName, payloadType, "", 1) + exportedFieldName + payloadType
			}

			fields = append(fields, exportedFieldName+" "+typeName+" "+tags)
			fieldStrings, nssFromRecursive := generateStructFields(payloadType, handlerName, v, []string{}, []myfile.NewStruct{}, typeName)
			nss = append(nss, myfile.NewStruct{
				Name:   typeName,
				Fields: fieldStrings,
			})
			nss = append(nss, nssFromRecursive...)
		case []interface{}:
			switch elem := v[0].(type) {
			case *mymap.OrderedMap:
				fieldReplace := strings.Replace(fieldName, "|", " ", 1)
				fieldParts := strings.Fields(fieldReplace)
				fieldName := fieldParts[0]
				required = len(fieldParts) > 1 && fieldParts[1] == "required"
				camelCaseFieldName := mystr.CamelCase(fieldName)
				exportedFieldName := mystr.CapitalizeFirstLetter(camelCaseFieldName)
				tags := "`json:\"" + fieldName + "\""
				if required {
					tags += " validate:\"required\""
				}
				tags += "`"

				// using previousTypeName to generate recursive struct name
				typeName := ""
				if previousTypeName == "" {
					typeName = "[]" + handlerName + exportedFieldName + payloadType
				} else {
					typeName = strings.Replace(previousTypeName, payloadType, "", 1) + exportedFieldName + payloadType
				}

				fields = append(fields, exportedFieldName+" "+typeName+" "+tags)
				fieldStrings, nssFromRecursive := generateStructFields(payloadType, handlerName, elem, []string{}, []myfile.NewStruct{}, "")
				nss = append(nss, myfile.NewStruct{
					Name:   strings.Replace(typeName, "[]", "", 1),
					Fields: fieldStrings,
				})
				nss = append(nss, nssFromRecursive...)
			case string:
				elemReplace := strings.Replace(elem, "|", " ", 1)
				typeParts := strings.Fields(elemReplace)
				baseType = "[]" + typeParts[0]
				required = len(typeParts) > 1 && typeParts[1] == "required"

				camelCaseFieldName := mystr.CamelCase(fieldName)
				exportedFieldName := mystr.CapitalizeFirstLetter(camelCaseFieldName)

				tags := "`json:\"" + fieldName + "\""
				if required {
					tags += " validate:\"required\""
				}
				tags += "`"

				fields = append(fields, exportedFieldName+" "+baseType+" "+tags)
			}
		}
	}
	return strings.Join(fields, "\n\t"), nss
}
