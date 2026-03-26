package testscommon

var (
	// TODO: use "testAccount", instead
	// TestAddressAlice is a test address
	TestAddressAlice = "drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf"
	// TestPubKeyAlice is a test pubkey
	TestPubKeyAlice, _ = RealWorldBech32PubkeyConverter.Decode(TestAddressAlice)

	// TODO: use "testAccount", instead
	// TestAddressBob is a test address
	TestAddressBob = "drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c"
	// TestPubKeyBob is a test pubkey
	TestPubKeyBob, _ = RealWorldBech32PubkeyConverter.Decode(TestAddressBob)

	// TODO: use "testAccount", instead
	// TestAddressCarol is a test address
	TestAddressCarol = "drt1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq889n6e"
	// TestPubKeyCarol is a test pubkey
	TestPubKeyCarol, _ = RealWorldBech32PubkeyConverter.Decode(TestAddressCarol)

	// TODO: use "testAccount", instead
	// TestAddressOfContract is a test address
	TestAddressOfContract = "drt1qqqqqqqqqqqqqpgqfejaxfh4ktp8mh8s77pl90dq0uzvh2vk396qzye6zs"

	// TestUserAShard0 is a test account (user)
	TestUserAShard0 = newTestAccount("drt1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqlqde3c")

	// TestUserBShard0 is a test account (user)
	TestUserBShard0 = newTestAccount("drt1uv40ahysflse896x4ktnh6ecx43u7cmy9wnxnvcyp7deg299a4sq8s28dr")

	// TestUserCShard0 is a test account (user)
	TestUserCShard0 = newTestAccount("drt1ncsyvhku3q7zy8f8rjmmx2t9zxgch38cel28kzg3m8pt86dt0vqqyyevt6")

	// TestContractFooShard0 is a test account (contract)
	TestContractFooShard0 = newTestAccount("drt1qqqqqqqqqqqqqpgqagjekf5mxv86hy5c62vvtug5vc6jmgcsq6uq6lwq7w")

	// TestContractBarShard0 is a test account (contract)
	TestContractBarShard0 = newTestAccount("drt1qqqqqqqqqqqqqpgqdstpe4fepzl4w8683xw88t5kcjkxz0zaq6uqpwdpgz")

	// TestContractFooShard1 is a test account (contract)
	TestContractFooShard1 = newTestAccount("drt1qqqqqqqqqqqqqpgq89t5xm4x04tnt9lv747wdrsaycf3rcwcggzsqz0qn8")

	// TestContractBarShard1 is a test account (contract)
	TestContractBarShard1 = newTestAccount("drt1qqqqqqqqqqqqqpgq0dtujxcrmwwqdtwzvq5nxuwgjcgaty7fggzsykmcc5")

	// TestContractFooShard2 is a test account (contract)
	TestContractFooShard2 = newTestAccount("drt1qqqqqqqqqqqqqpgqeesfamasje5zru7ku79m8p4xqfqnywvqxj0q2hnpwa")

	// TestContractBarShard2 is a test account (contract)
	TestContractBarShard2 = newTestAccount("drt1ux2wvqyh8pw8ea26urjqq65mqytfn42dr980pvucztxk9w79xj0qh6jyg8")

	// TestUserShard1 is a test account (user)
	TestUserShard1 = newTestAccount("drt1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssey5egf")

	// TestUserShard2 is a test account (user)
	TestUserShard2 = newTestAccount("drt1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq889n6e")
)

type testAccount struct {
	Address string
	PubKey  []byte
}

func newTestAccount(address string) *testAccount {
	pubKey, _ := RealWorldBech32PubkeyConverter.Decode(address)

	return &testAccount{
		Address: address,
		PubKey:  pubKey,
	}
}
