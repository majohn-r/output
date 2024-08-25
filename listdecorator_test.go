package output

import "testing"

func Test_listDecorator_decorator(t *testing.T) {
	type fields struct {
		enabled    bool
		numeric    bool
		itemNumber uint8
	}
	tests := map[string]struct {
		fields         fields
		want           string
		wantItemNumber uint8
	}{
		"disabled": {
			fields: fields{
				enabled:    false,
				numeric:    false,
				itemNumber: 0,
			},
			want:           "",
			wantItemNumber: 0,
		},
		"enabled bulleted": {
			fields: fields{
				enabled:    true,
				numeric:    false,
				itemNumber: 0,
			},
			want:           "● ",
			wantItemNumber: 0,
		},
		"enabled numeric": {
			fields: fields{
				enabled:    true,
				numeric:    true,
				itemNumber: 3,
			},
			want:           " 3. ",
			wantItemNumber: 4,
		},
	}
	for name, tt := range tests {
		t.Run(name, func(t *testing.T) {
			ld := newListDecorator(tt.fields.enabled, tt.fields.numeric)
			ld.itemNumber = tt.fields.itemNumber
			if got := ld.Decorator(); got != tt.want {
				t.Errorf("Decorator() = %v, want %v", got, tt.want)
			}
			if got := ld.itemNumber; got != tt.wantItemNumber {
				t.Errorf("itemNumber = %v, want %v", got, tt.wantItemNumber)
			}
		})
	}
}
