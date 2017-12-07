package schema

//
//import (
//	"github.com/graphql-go/graphql"
//	"github.com/mappcpd/web-services/cmd/webd/graphql/data"
//	"github.com/mappcpd/web-services/cmd/webd/graphql/schema/types"
//)
//
//var Members = &graphql.Field{
//	Name:        "Members",
//	Description: "Fetch a list of members",
//	Type:        graphql.NewList(Member),
//	// todo - implement search arg
//	//Args: graphql.FieldConfigArgument{
//	//	"id": &graphql.ArgumentConfig{
//	//		Type:        &graphql.NonNull{OfType: graphql.Int},
//	//		Description: "Vessel UUID",
//	//	},
//	//},
//	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
//		//id, ok := p.Args["id"].(string)
//		//if ok {
//		//	return data.GetVessel(id), nil
//		//}
//		//return nil, nil
//		//xm := data.GetMembers()
//		//query := map[string]interface{}{
//		//	"lastName": "Smith",
//		//}
//		//proj := map[string]interface{}{}
//		//lim := 20
//		xm := data.GetMembers()
//		//if err != nil {
//		//	fmt.Println(err)
//		//}
//
//		return xm, nil
//	},
//}
