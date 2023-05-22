package resources

import (
	"fmt"
	"os"
)

type(
	Resource string)

func ReadResources(path string)([]Resource,error){
	f,err := os.OpenFile(path, os.O_RDONLY, os.ModeExclusive)
	if err != nil{
		return nil,err
	}

	defer f.Close()

	var res []Resource
	var val string
	for{
		if _,err := fmt.Fscanln(f, &val); err != nil{
			break
		}
		
		res = append(res, Resource(val))
	}

	return res, nil
}