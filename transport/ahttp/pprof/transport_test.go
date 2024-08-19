package pprof

import (
	"testing"
)

var testStr = `
<html>
<head>
<title>/debug/pprof/</title>
<style>
.profile-name{
	display:inline-block;
	width:6rem;
}
</style>
</head>
<body>
/debug/pprof/
<br>
<p>Set debug=1 as a query parameter to export in legacy text format</p>
<br>
Types of profiles available:
<table>
<thead><td>Count</td><td>Profile</td></thead>
<tr><td>317</td><td><a href='allocs?debug=1'>allocs</a></td></tr>
<tr><td>0</td><td><a href='block?debug=1'>block</a></td></tr>
<tr><td>0</td><td><a href='cmdline?debug=1'>cmdline</a></td></tr>
<tr><td>569</td><td><a href='goroutine?debug=1'>goroutine</a></td></tr>
<tr><td>317</td><td><a href='heap?debug=1'>heap</a></td></tr>
<tr><td>0</td><td><a href='mutex?debug=1'>mutex</a></td></tr>
<tr><td>0</td><td><a href='profile?debug=1'>profile</a></td></tr>
<tr><td>10</td><td><a href='threadcreate?debug=1'>threadcreate</a></td></tr>
<tr><td>0</td><td><a href='trace?debug=1'>trace</a></td></tr>
</table>
<a href="goroutine?debug=2">full goroutine stack dump</a>
<br>
<p>
Profile Descriptions:
<ul>
<li><div class=profile-name>allocs: </div> A sampling of all past memory allocations</li>
<li><div class=profile-name>block: </div> Stack traces that led to blocking on synchronization primitives</li>
<li><div class=profile-name>cmdline: </div> The command line invocation of the current program</li>
<li><div class=profile-name>goroutine: </div> Stack traces of all current goroutines. Use debug=2 as a query parameter to export in the same format as an unrecovered panic.</li>
<li><div class=profile-name>heap: </div> A sampling of memory allocations of live objects. You can specify the gc GET parameter to run GC before taking the heap sample.</li>
<li><div class=profile-name>mutex: </div> Stack traces of holders of contended mutexes</li>
<li><div class=profile-name>profile: </div> CPU profile. You can specify the duration in the seconds GET parameter. After you get the profile file, use the go tool pprof command to investigate the profile.</li>
<li><div class=profile-name>threadcreate: </div> Stack traces that led to the creation of new OS threads</li>
<li><div class=profile-name>trace: </div> A trace of execution of the current program. You can specify the duration in the seconds GET parameter. After you get the trace file, use the go tool trace command to investigate the trace.</li>
</ul>
</p>
</body>
</html>
`

func Test_transporter_rebuildHref(t *testing.T) {
	// reg := regexp.MustCompile(`href='([\w?=]*)'`)
	// result := reg.FindAllStringSubmatch(testStr, -1)

	// str := testStr
	// for _, val := range result {
	// 	newStr := val[1] + "&ID=101"
	// 	str = strings.Replace(str, val[1], newStr, -1)
	// }
	// t.Logf("%s", cast.ToJson(result))
	// t.Logf("----------------------")

	trans := NewTransporter("", "/server_test", WithParam("ID", "1001"))
	str := trans.rebuildMainPage([]byte(testStr))
	t.Logf("%s", str)
}

func Test_transporter_getPprofType(t *testing.T) {
	trans := NewTransporter("", "/server_test", WithParam("ID", "1001"))

	type args struct {
		rawURL string
	}
	tests := []struct {
		name  string
		trans *Transporter
		args  args
		want  string
	}{
		{
			name:  "test1",
			trans: trans,
			args: args{
				rawURL: "heap?debug=1",
			},
			want: "heap",
		},
		{
			name:  "test2",
			trans: trans,
			args: args{
				rawURL: "goroutine?debug=1&ID=1001",
			},
			want: "goroutine",
		},
		{
			name:  "test3",
			trans: trans,
			args: args{
				rawURL: "cmdline",
			},
			want: "cmdline",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.trans.getPprofType(tt.args.rawURL); got != tt.want {
				t.Errorf("transporter.getPprofType() = %v, want %v", got, tt.want)
			}
		})
	}
}
