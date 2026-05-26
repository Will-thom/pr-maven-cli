package prmaven

import "testing"

func TestSlashPathNormalizesWindowsAndPOSIXSeparators(t *testing.T) {
	tests := []struct {
		name string
		path string
		want string
	}{
		{
			name: "windows style relative path",
			path: `platform\service-core\target\surefire-reports\TEST-dev.prmaven.demo.NestedPaymentTest.xml`,
			want: "platform/service-core/target/surefire-reports/TEST-dev.prmaven.demo.NestedPaymentTest.xml",
		},
		{
			name: "posix style relative path",
			path: "platform/service-core/target/surefire-reports/TEST-dev.prmaven.demo.NestedPaymentTest.xml",
			want: "platform/service-core/target/surefire-reports/TEST-dev.prmaven.demo.NestedPaymentTest.xml",
		},
		{
			name: "mixed separators",
			path: `platform/service-core\target/surefire-reports`,
			want: "platform/service-core/target/surefire-reports",
		},
		{
			name: "current directory",
			path: ".",
			want: ".",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := slashPath(tt.path)
			if got != tt.want {
				t.Fatalf("slashPath(%q) = %q, want %q", tt.path, got, tt.want)
			}
		})
	}
}
