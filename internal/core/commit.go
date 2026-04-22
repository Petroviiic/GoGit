package core

type Commit struct {
	BaseObject
	ParentHash string // Dodatno polje koje Blob nema
	Author     string
}

// func NewObject(objectType string, content []byte) *Object {
// 	return &Object{

// 		deserialize: func(compressedData []byte) {
// 			b := bytes.NewReader(compressedData)

// 			r, err := zlib.NewReader(b)
// 			if err != nil {
// 				panic(err)
// 			}
// 			defer r.Close()

// 			io.Copy(os.Stdout, r)
// 		},
// 	}
// }
