package modes

import "testing"

func TestWORDPRESS(t *testing.T) {
	type args struct {
		hash  string
		plain string
		in2   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"$P$946647711V1klyitUYhtB8Yw5DMA/w.", "hashcat", ""}, true},
		{"", args{"$P$946647711V1klyitUYhtB8Yw5DMA/w.", "hashcat2", ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := WORDPRESS(tt.args.hash, tt.args.plain, tt.args.in2); got != tt.want {
				t.Errorf("WORDPRESS() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMD5(t *testing.T) {
	type args struct {
		h     string
		plain string
		in2   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"8743b52063cd84097a65d1633f5c74f5", "hashcat", ""}, true},
		{"", args{"8743b52063cd84097a65d1633f5c74f5", "hashcat2", ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5(tt.args.h, tt.args.plain, tt.args.in2); got != tt.want {
				t.Errorf("MD5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMD5PLAINSALT(t *testing.T) {
	type args struct {
		h     string
		plain string
		salt  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"3d83c8e717ff0e7ecfe187f088d69954", "hashcat", "343141"}, true},
		{"", args{"3d83c8e717ff0e7ecfe187f088d69954", "hashcat2", "343141"}, false},

		// JOOMLA < 2.5.18
		{"", args{"b78f863f2c67410c41e617f724e22f34", "hashcat", "89384528665349271307465505333378"}, true},
		{"", args{"b78f863f2c67410c41e617f724e22f34", "hashcat2", "89384528665349271307465505333378"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5PLAINSALT(tt.args.h, tt.args.plain, tt.args.salt); got != tt.want {
				t.Errorf("MD5PLAINSALT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMD5SALTPLAIN(t *testing.T) {
	type args struct {
		h     string
		plain string
		salt  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"57ab8499d08c59a7211c77f557bf9425", "hashcat", "4247"}, true},
		{"", args{"57ab8499d08c59a7211c77f557bf9425", "hashcat2", "4247"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5SALTPLAIN(tt.args.h, tt.args.plain, tt.args.salt); got != tt.want {
				t.Errorf("MD5SALTPLAIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA1(t *testing.T) {
	type args struct {
		hash  string
		plain string
		in2   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"b89eaac7e61417341b710b727768294d0e6a277b", "hashcat", ""}, true},
		{"", args{"b89eaac7e61417341b710b727768294d0e6a277b", "hashcat2", ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SHA1(tt.args.hash, tt.args.plain, tt.args.in2); got != tt.want {
				t.Errorf("SHA1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA1SALTPLAIN(t *testing.T) {
	type args struct {
		hash  string
		plain string
		salt  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"a428863972744b16afef28e0087fc094b44bb7b1", "hashcat", "465727565"}, true},
		{"", args{"a428863972744b16afef28e0087fc094b44bb7b1", "hashcat2", "465727565"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SHA1SALTPLAIN(tt.args.hash, tt.args.plain, tt.args.salt); got != tt.want {
				t.Errorf("SHA1SALTPLAIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestVBULLETIN(t *testing.T) {
	type args struct {
		h     string
		plain string
		salt  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"2611", args{"28f9975808ae2bdc5847b1cda26033ea", "hashcat", "308"}, true},
		{"2611", args{"28f9975808ae2bdc5847b1cda26033ea", "hashcat2", "308"}, false},

		{"2711", args{"0844fbb2fdeda31884a7a45ec2010bb6", "hashcat", "324410183853308365427804872426"}, true},
		{"2711", args{"0844fbb2fdeda31884a7a45ec2010bb6", "hashcat2", "324410183853308365427804872426"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := VBULLETIN(tt.args.h, tt.args.plain, tt.args.salt); got != tt.want {
				t.Errorf("VBULLETIN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestIPB(t *testing.T) {
	type args struct {
		h     string
		plain string
		salt  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"2811", args{"022f7e02b3314f7d0968f73c00ba759f", "hashcat", "67588"}, true},
		{"2811", args{"022f7e02b3314f7d0968f73c00ba759f", "hashcat2", "67588"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IPB(tt.args.h, tt.args.plain, tt.args.salt); got != tt.want {
				t.Errorf("IPB() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestARGON2(t *testing.T) {
	type args struct {
		hash  string
		plain string
		in2   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"argon2i positive", args{"$argon2i$v=19$m=16,t=2,p=1$eHYyOFNndmFhd0NJZjZ4cw$S79XfFyFqTSlZsZe9338JQ", "hashkitty", ""}, true},
		{"argon2i negative", args{"$argon2i$v=19$m=16,t=2,p=1$eHYyOFNndmFhd0NJZjZ4cw$S79XfFyFqTSlZsZe9338JQ", "hashkitty2", ""}, false},

		// TODO: argon2d seems to be broken atm
		//{"argon2d positive", args{"$argon2d$v=19$m=16,t=2,p=1$cWxIVTE3TE1BSWpsMTNJQw$VH2sp49ndEqUYUaz/HhgIw", "hashkitty", ""}, true},
		//{"argon2d negative", args{"$argon2d$v=19$m=16,t=2,p=1$cWxIVTE3TE1BSWpsMTNJQw$VH2sp49ndEqUYUaz/HhgIw", "hashkitty2", ""}, false},

		{"argon2id positive", args{"$argon2id$v=19$m=16,t=2,p=1$eHYyOFNndmFhd0NJZjZ4cw$qpk8FFWSTXf5EV6ndCvRpg", "hashkitty", ""}, true},
		{"argon2id negative", args{"$argon2id$v=19$m=16,t=2,p=1$eHYyOFNndmFhd0NJZjZ4cw$qpk8FFWSTXf5EV6ndCvRpg", "hashkitty2", ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := ARGON2(tt.args.hash, tt.args.plain, tt.args.in2); got != tt.want {
				t.Errorf("ARGON2() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA1DASH(t *testing.T) {
	type args struct {
		hash  string
		plain string
		salt  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"d9e8b6c25ee0760f1b90ffa14eafbe09aede0e65", "hashkitty", "tlas"}, true},
		{"", args{"d9e8b6c25ee0760f1b90ffa14eafbe09aede0e65", "hashkitty2", "tlas"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SHA1DASH(tt.args.hash, tt.args.plain, tt.args.salt); got != tt.want {
				t.Errorf("SHA1DASH() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBCRYPT(t *testing.T) {
	type args struct {
		hash  string
		plain string
		in2   string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"$2a$12$jbrdD9p.tKWyqXZSFCqiS.LoohuHrkM2fRD/MCD6IDFag2ciD.eFa", "hashkitty", ""}, true},
		{"", args{"$2a$12$jbrdD9p.tKWyqXZSFCqiS.LoohuHrkM2fRD/MCD6IDFag2ciD.eFa", "hashkitty2", ""}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := BCRYPT(tt.args.hash, tt.args.plain, tt.args.in2); got != tt.want {
				t.Errorf("BCRYPT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA256PLAINSALT(t *testing.T) {
	type args struct {
		h     string
		plain string
		salt  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"04ca9386cffb05ec8ec9ba492ef73bf3e86f445e898d9b33b9a56ebcfe7c5463", "hashkitty", "tlas"}, true},
		{"", args{"04ca9386cffb05ec8ec9ba492ef73bf3e86f445e898d9b33b9a56ebcfe7c5463", "hashkitty2", "tlas"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SHA256PLAINSALT(tt.args.h, tt.args.plain, tt.args.salt); got != tt.want {
				t.Errorf("SHA256PLAINSALT() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMYBB(t *testing.T) {
	type args struct {
		h     string
		plain string
		salt  string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{"", args{"cd74cfc96088f9d4fccdf9bdc16cb3ec", "hashkitty", "tlas"}, true},
		{"", args{"cd74cfc96088f9d4fccdf9bdc16cb3ec", "hashkitty2", "tlas"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MYBB(tt.args.h, tt.args.plain, tt.args.salt); got != tt.want {
				t.Errorf("MYBB() = %v, want %v", got, tt.want)
			}
		})
	}
}
