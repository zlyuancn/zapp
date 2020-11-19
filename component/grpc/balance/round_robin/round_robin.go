/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2020/11/19
   Description :
-------------------------------------------------
*/

package round_robin

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer/roundrobin"
)

const Name = roundrobin.Name

func Balance() grpc.DialOption {
	return grpc.WithDefaultServiceConfig(fmt.Sprintf(`{ "loadBalancingConfig": [{"%v": {}}] }`, roundrobin.Name))
}
