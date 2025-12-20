package utils

import (
	"testing"
)

func TestIsGitHubRawURL(t *testing.T) {
	tests := []struct {
		name string
		url  string
		want bool
	}{
		{
			name: "HTTPS GitHub raw URL",
			url:  "https://raw.githubusercontent.com/wangtao53-tal/image/main/tango/IMG_5829.JPG",
			want: true,
		},
		{
			name: "HTTP GitHub raw URL",
			url:  "http://raw.githubusercontent.com/user/repo/branch/path.jpg",
			want: true,
		},
		{
			name: "Non-GitHub URL",
			url:  "https://example.com/image.jpg",
			want: false,
		},
		{
			name: "Data URL",
			url:  "data:image/jpeg;base64,/9j/4AAQSkZJRg==",
			want: false,
		},
		{
			name: "Empty string",
			url:  "",
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsGitHubRawURL(tt.url); got != tt.want {
				t.Errorf("IsGitHubRawURL() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestConvertToJSDelivrCDN(t *testing.T) {
	tests := []struct {
		name    string
		rawURL  string
		wantCDN string
		wantErr bool
	}{
		{
			name:    "Valid GitHub raw URL",
			rawURL:  "https://raw.githubusercontent.com/wangtao53-tal/image/main/tango/IMG_5829.JPG",
			wantCDN: "https://cdn.jsdelivr.net/gh/wangtao53-tal/image@main/tango/IMG_5829.JPG",
			wantErr: false,
		},
		{
			name:    "GitHub raw URL with subdirectory",
			rawURL:  "https://raw.githubusercontent.com/user/repo/branch/path/to/image.png",
			wantCDN: "https://cdn.jsdelivr.net/gh/user/repo@branch/path/to/image.png",
			wantErr: false,
		},
		{
			name:    "Non-GitHub URL",
			rawURL:  "https://example.com/image.jpg",
			wantCDN: "https://example.com/image.jpg",
			wantErr: false,
		},
		{
			name:    "Invalid GitHub raw URL format",
			rawURL:  "https://raw.githubusercontent.com/user",
			wantCDN: "https://raw.githubusercontent.com/user",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ConvertToJSDelivrCDN(tt.rawURL)
			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertToJSDelivrCDN() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.wantCDN {
				t.Errorf("ConvertToJSDelivrCDN() = %v, want %v", got, tt.wantCDN)
			}
		})
	}
}

func TestExtractGitHubRawURLInfo(t *testing.T) {
	tests := []struct {
		name       string
		rawURL     string
		wantOwner  string
		wantRepo   string
		wantBranch string
		wantPath   string
		wantOk     bool
	}{
		{
			name:       "Valid GitHub raw URL",
			rawURL:     "https://raw.githubusercontent.com/wangtao53-tal/image/main/tango/IMG_5829.JPG",
			wantOwner:  "wangtao53-tal",
			wantRepo:   "image",
			wantBranch: "main",
			wantPath:   "tango/IMG_5829.JPG",
			wantOk:     true,
		},
		{
			name:       "GitHub raw URL with HTTP",
			rawURL:     "http://raw.githubusercontent.com/user/repo/branch/path.jpg",
			wantOwner:  "user",
			wantRepo:   "repo",
			wantBranch: "branch",
			wantPath:   "path.jpg",
			wantOk:     true,
		},
		{
			name:       "Non-GitHub URL",
			rawURL:     "https://example.com/image.jpg",
			wantOwner:  "",
			wantRepo:   "",
			wantBranch: "",
			wantPath:   "",
			wantOk:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotOwner, gotRepo, gotBranch, gotPath, gotOk := ExtractGitHubRawURLInfo(tt.rawURL)
			if gotOwner != tt.wantOwner {
				t.Errorf("ExtractGitHubRawURLInfo() gotOwner = %v, want %v", gotOwner, tt.wantOwner)
			}
			if gotRepo != tt.wantRepo {
				t.Errorf("ExtractGitHubRawURLInfo() gotRepo = %v, want %v", gotRepo, tt.wantRepo)
			}
			if gotBranch != tt.wantBranch {
				t.Errorf("ExtractGitHubRawURLInfo() gotBranch = %v, want %v", gotBranch, tt.wantBranch)
			}
			if gotPath != tt.wantPath {
				t.Errorf("ExtractGitHubRawURLInfo() gotPath = %v, want %v", gotPath, tt.wantPath)
			}
			if gotOk != tt.wantOk {
				t.Errorf("ExtractGitHubRawURLInfo() gotOk = %v, want %v", gotOk, tt.wantOk)
			}
		})
	}
}
