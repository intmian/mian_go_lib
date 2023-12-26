package pushmod

import (
	"testing"
)

func TestMgr_PushPushDeer(t *testing.T) {
	type fields struct {
		pushEmailToken *EmailSetting
		PushDeerToken  *PushDeerSetting
	}
	type args struct {
		title    string
		content  string
		markDown bool
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   error
	}{
		{
			name: "public-text",
			fields: fields{
				pushEmailToken: nil,
				PushDeerToken: &PushDeerSetting{
					Token: "PDU10120Tp8PByEPFdrKiStSvMWeOdeFtwY7GuOmQ",
				},
			},
			args: args{
				title:    "test",
				content:  "contentTest",
				markDown: false,
			},
			want: nil,
		},
		{
			name: "public-markdown",
			fields: fields{
				pushEmailToken: nil,
				PushDeerToken: &PushDeerSetting{
					Token: "PDU10120Tp8PByEPFdrKiStSvMWeOdeFtwY7GuOmQ",
				},
			},
			args: args{
				title:    "test",
				content:  "# 论信息\n\n- 很重要\n- 测试aaa111\n\n> 兼听则明\n\n```c++\n偏信则黯\n```\n",
				markDown: true,
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m, err := NewPushDeerMgr(&PushDeerSetting{
				Token: tt.fields.PushDeerToken.Token,
			})
			if err != nil {
				t.Fatal(err)
			}
			if tt.args.markDown {
				err = m.PushMarkDown(tt.args.title, tt.args.content)
			} else {
				err = m.Push(tt.args.title, tt.args.content)
			}
			if err != nil {
				t.Fatal(err)
			}
		})
	}
}
