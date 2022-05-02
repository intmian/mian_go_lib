package push

import "testing"

func TestMgr_PushPushDeer(t *testing.T) {
	type fields struct {
		pushEmailToken *EmailToken
		PushDeerToken  *PushDeerToken
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
		want   string
		want1  bool
	}{
		{
			name: "public-text",
			fields: fields{
				pushEmailToken: nil,
				PushDeerToken: &PushDeerToken{
					Token: "PDU10120Tp8PByEPFdrKiStSvMWeOdeFtwY7GuOmQ",
				},
			},
			args: args{
				title:    "test",
				content:  "contentTest",
				markDown: false,
			},
			want:  "200 OK",
			want1: true,
		},
		{
			name: "public-markdown",
			fields: fields{
				pushEmailToken: nil,
				PushDeerToken: &PushDeerToken{
					Token: "PDU10120Tp8PByEPFdrKiStSvMWeOdeFtwY7GuOmQ",
				},
			},
			args: args{
				title:    "test",
				content:  "# 论信息\n\n- 很重要\n- 测试aaa111\n\n> 兼听则明\n\n```c++\n偏信则黯\n```\n",
				markDown: true,
			},
			want:  "200 OK",
			want1: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &Mgr{
				pushEmailToken: tt.fields.pushEmailToken,
				pushDeerToken:  tt.fields.PushDeerToken,
			}
			got, got1 := m.PushPushDeer(tt.args.title, tt.args.content, tt.args.markDown)
			if got != tt.want {
				t.Errorf("PushPushDeer() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("PushPushDeer() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
