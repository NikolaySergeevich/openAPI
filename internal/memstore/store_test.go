package memstore

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

func TestStore_Add(t *testing.T){
	t.Parallel()
	expLen := 1
	exp := Item{
		ID: uuid.New().String(),
		Name: "test_name",
		Lat: 10.10,
		Lon: 34.44,

	}

	s := New()
	s.Add(exp)
	var found bool
	for _, v := range s.items {
		if v.ID == exp.ID {
			found = true
			break
		}
	}

	
	//а теперь применение testify:
	
	assert.Lenf(t, s.items, expLen, "длина items меньше чем exp")
	
	assert.True(t, found, "элементы не совпали")//эта и след. строка несут одну и ту же логику.
	assert.Equal(t, found, true, "element not found")

	if !found {
		t.Errorf("got: %T, want true", found)
		return
	}
}

func TestStoreFindAll (t *testing.T){
	t.Parallel()
//смысл теста в том, что бы проверить правильно ли вернёт функция FindAll значения при разном кол-ве объекте
//для этого создаём структуру для разной ситуации: когда только один объект
//когда много объектов, когда nil и когда просто пусто.
	testCases := []struct{ 
		name string
		items []Item
		expLen int
	}{
		{
			name: "test_singl_item",
			items: []Item{
				{
					ID: uuid.New().String(),
					Name: "test 1",
				},
			},
			expLen: 1,
		},
		{
			name: "test_many_item",
			items: []Item{
				{
					ID: uuid.New().String(),
					Name: "test 1",
				},
				{
					ID: uuid.New().String(),
					Name: "test 2",
				},
				{
					ID: uuid.New().String(),
					Name: "test 3",
				},
			},
			expLen: 3,
		},
		{
			name: "test_nil_items",
			items: nil,
			expLen: 0,
		},
		{
			name: "test_emty_items",
			items: make([]Item, 0),
			expLen: 0,
		},
	}
	//теперь итерируясь по этой структуре, будет применяться 
	//метод FindAll для каждого случая.
	for _, tc := range testCases {
		tc := tc//мы используем переменность, поэтому нужно пересоздавать переменную, т.к вызываться метод бужет в замыкании
		t.Run(//при итерации будет запускаться паралельно отдельно тестик для каждой ситуации
			tc.name, func(t *testing.T){
				t.Parallel()
				s := New()//создаётся структура 
				s.items = append(s.items, tc.items...)//в структуру кладётся набор items, котрые создавались при итерации
				assert.Lenf(t, s.FindAll(), tc.expLen, "Функция FindAll отработала не так как нужно. Должна была найти %d элементов, а нашла только %d", tc.expLen, len(s.FindAll()))
				//тут может быть вопрос. s.FindAll() вернёт слайс, а tc.expLen это просто int, как тогда функция assert.Len сравнит это?
				//эта функция возьмёт длину s.FindAll() и сравнит с tc.expLen
			},
		)
	}

	

}