package misc

import "testing"

func TestGetIndexPage(t *testing.T) {
	type args struct {
		index      int
		pageSize   int
		pageBegin1 bool
	}
	tests := []struct {
		name          string
		args          args
		wantPage      int
		wantPageIndex int
	}{
		{
			name: "1",
			args: args{
				index:      21,
				pageSize:   10,
				pageBegin1: true,
			},
			wantPage:      3,
			wantPageIndex: 1,
		},
		{
			name: "2",
			args: args{
				index:      21,
				pageSize:   10,
				pageBegin1: false,
			},
			wantPage:      2,
			wantPageIndex: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPage, gotPageIndex := GetIndexPage(tt.args.index, tt.args.pageSize, tt.args.pageBegin1)
			if gotPage != tt.wantPage {
				t.Errorf("GetIndexPage() gotPage = %v, want %v", gotPage, tt.wantPage)
			}
			if gotPageIndex != tt.wantPageIndex {
				t.Errorf("GetIndexPage() gotPageIndex = %v, want %v", gotPageIndex, tt.wantPageIndex)
			}
		})
	}
}

func TestGetOriIndex(t *testing.T) {
	type args struct {
		index      int
		page       int
		pageSize   int
		pageBegin1 bool
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{
			name: "1",
			args: args{
				index:      2,
				page:       2,
				pageSize:   10,
				pageBegin1: false,
			},
			want: 22,
		},
		{
			name: "2",
			args: args{
				index:      2,
				page:       2,
				pageSize:   10,
				pageBegin1: true,
			},
			want: 11,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := GetPageIndexOriIndex(tt.args.index, tt.args.page, tt.args.pageSize, tt.args.pageBegin1); got != tt.want {
				t.Errorf("GetPageIndexOriIndex() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPageStartEnd(t *testing.T) {
	type args struct {
		page       int
		pageSize   int
		total      int
		pageBegin1 bool
	}
	tests := []struct {
		name          string
		args          args
		wantPageStart int
		wantPageEnd   int
	}{
		{
			name: "1",
			args: args{
				page:       1,
				pageSize:   10,
				total:      15,
				pageBegin1: true,
			},
			wantPageStart: 1,
			wantPageEnd:   10,
		},
		{
			name: "2",
			args: args{
				page:       2,
				pageSize:   10,
				total:      15,
				pageBegin1: true,
			},
			wantPageStart: 11,
			wantPageEnd:   15,
		},
		{
			name: "3",
			args: args{
				page:       0,
				pageSize:   10,
				total:      15,
				pageBegin1: false,
			},
			wantPageStart: 0,
			wantPageEnd:   9,
		},
		{
			name: "4",
			args: args{
				page:       1,
				pageSize:   10,
				total:      15,
				pageBegin1: false,
			},
			wantPageStart: 10,
			wantPageEnd:   14,
		},
		{
			name: "5",
			args: args{
				page:       2,
				pageSize:   10,
				total:      15,
				pageBegin1: false,
			},
			wantPageStart: -1,
			wantPageEnd:   -1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPageStart, gotPageEnd := GetPageStartEnd(tt.args.page, tt.args.pageSize, tt.args.total, tt.args.pageBegin1)
			if gotPageStart != tt.wantPageStart {
				t.Errorf("GetPageStartEnd() gotPageStart = %v, want %v", gotPageStart, tt.wantPageStart)
			}
			if gotPageEnd != tt.wantPageEnd {
				t.Errorf("GetPageStartEnd() gotPageEnd = %v, want %v", gotPageEnd, tt.wantPageEnd)
			}
		})
	}
}

func TestGetMaxPage(t *testing.T) {
	type args struct {
		total      int
		pageSize   int
		pageBegin1 bool
	}
	tests := []struct {
		name        string
		args        args
		wantMaxPage int
	}{
		{
			name: "1",
			args: args{
				total:      15,
				pageSize:   10,
				pageBegin1: true,
			},
			wantMaxPage: 2,
		},
		{
			name: "2",
			args: args{
				total:      15,
				pageSize:   10,
				pageBegin1: false,
			},
			wantMaxPage: 1,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if gotMaxPage := GetMaxPage(tt.args.total, tt.args.pageSize, tt.args.pageBegin1); gotMaxPage != tt.wantMaxPage {
				t.Errorf("GetMaxPage() = %v, want %v", gotMaxPage, tt.wantMaxPage)
			}
		})
	}
}
