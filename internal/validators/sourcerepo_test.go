package validators

//func TestSourceRepoValidator_Validate(t *testing.T) {
//	var cases = []struct {
//		in   string
//		repo string
//		out  bool
//	}{
//		{
//			"./testdata/valid.slsa.json",
//			"pkg:github/foo/bar@deadbeef",
//			true,
//		},
//		{
//			"./testdata/small_gitlab.slsa.json",
//			"pkg:gitlab/foo%2Fbar/zoo@deadbeef",
//			true,
//		},
//		{
//			"./testdata/invalid_sourcerepo.slsa.json",
//			"pkg:github/foo/bar@deadbeef",
//			false,
//		},
//		{
//			"./testdata/valid.slsa.json",
//			"https://github.com/foo/bar.git@deadbeef",
//			true,
//		},
//		{
//			"./testdata/small_gitlab.slsa.json",
//			"https://gitlab.com/foo/bar/zoo.git@deadbeef",
//			true,
//		},
//		{
//			"./testdata/invalid_sourcerepo.slsa.json",
//			"https://github.com/foo/bar.git@deadbeef",
//			false,
//		},
//	}
//	for _, tt := range cases {
//		t.Run(tt.in, func(t *testing.T) {
//			v := SourceRepoValidator{Expected: tt.repo}
//			ok := v.Check1(loadFile(t, tt.in))
//			assert.EqualValues(t, tt.out, ok)
//		})
//	}
//}
