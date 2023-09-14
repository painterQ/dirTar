func TestTar(t *testing.T) {
	defer func() {
		_ = os.RemoveAll("./test")
	}()
	assert.Nil(t, EncryptDir(key, ".git", "./test/target.tar"))
	assert.Nil(t, DeTarDir(key, "./test/target.tar", "./test"))
}
