package pk

type Cipher interface {

	//Init generates and saves keys both private
	//and public to pem files for future use
	//Init(credentialsDir string) (err error)

	//Encode takes password provided by user as input
	//encode it abd return bytes that are stored
	//to the db for later retrieval
	Encode(password string) ([]byte, error)

	//Decode takes persisted data from database and decode it
	//to recover passwords that were stored
	Decode(encoded []byte) (string, error)

	//Sign generate message digest from password and sign it
	Sign(password string) (digest []byte, signature []byte, err error)

	//Verify takes password from Decode and generate a digest from it
	//then compares the digest with the one in db, if they match it
	//verify the signature of the current digest with the one in db
	//and return a non-nl error if they match
	Verify(password string, digest []byte, signature []byte) error
}


