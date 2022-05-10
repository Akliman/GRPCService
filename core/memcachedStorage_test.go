package core

import (
	"GRPCService/external/memcache"
	"context"
	"strconv"
	"sync"
	"testing"
)

//Паралельный запуск тестов в этом файле не поддерживается (
//Для выполнения тестов нужно запустить memcached сервер

const memcahedAdr = "localhost:11211" //Адрес для сервера memcached

func TestMemcachedStorage_Delete(t *testing.T) {
	err := CreateMemcahed(context.Background(), memcahedAdr)
	if err != nil {
		t.Error(err)
	}

	data := map[int32]string{
		1:    "ts1",
		2:    "ts2",
		3:    "ts3",
		4:    "ts4",
		5:    "ts5",
		100:  "ts100",
		1000: "ts1000",
	}

	//Очищаем хранилище
	err = memcache.FlushAll()
	if err != nil {
		t.Error(err)
	}

	for key, val := range data {
		err := memcache.Set(strconv.Itoa(int(key)), val, 0)
		if err != nil {
			t.Error(err)
		}
	}

	tests := []struct {
		name    string
		key     int32
		wantErr bool
	}{
		{"1",
			1,
			false},
		{"2",
			2,
			false},
		{"3",
			3,
			false},
		{"4",
			100,
			false},
		{"5",
			100,
			true},
		{"6",
			10000,
			true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MemcachedCurrentStorage.Delete(context.Background(), tt.key); (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	notDeletedKeys := []int32{4, 5, 1000}
	for _, key := range notDeletedKeys {
		val, err := memcache.Get(strconv.Itoa(int(key)))
		if err != nil {
			t.Error(err)
		}

		if val != "ts"+strconv.Itoa(int(key)) {
			t.Error("value is wrong")
		}
	}

}

func TestMemcachedStorage_Get(t *testing.T) {
	err := CreateMemcahed(context.Background(), memcahedAdr)
	if err != nil {
		t.Error(err)
	}

	data := map[int32]string{
		1:    "ts1",
		2:    "ts2",
		3:    "ts3",
		4:    "ts4",
		5:    "ts5",
		100:  "ts100",
		1000: "ts1000",
	}

	//Очищаем хранилище
	err = memcache.FlushAll()
	if err != nil {
		t.Error(err)
	}

	for key, val := range data {
		err := memcache.Set(strconv.Itoa(int(key)), val, 0)
		if err != nil {
			t.Error(err)
		}
	}

	tests := []struct {
		name    string
		key     int32
		want    string
		wantErr bool
	}{{"1",
		1,
		"ts1",
		false},
		{
			"2",
			2,
			"ts2",
			false,
		},
		{
			"3",
			1000,
			"ts1000",
			false,
		},
		{
			"4",
			6,
			"",
			true,
		},
		{
			"5",
			1001,
			"",
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := MemcachedCurrentStorage.Get(context.Background(), tt.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Get() got = %v, want %v", got, tt.want)
			}
		})
	}

	keys := []int32{1, 2, 3, 4, 5, 100, 1000}
	for _, key := range keys {
		val, err := memcache.Get(strconv.Itoa(int(key)))
		if err != nil {
			t.Error(err)
		}

		if val != "ts"+strconv.Itoa(int(key)) {
			t.Error("value is wrong")
		}
	}
}

func TestMemcachedStorage_Put(t *testing.T) {
	err := CreateMemcahed(context.Background(), memcahedAdr)
	if err != nil {
		t.Error(err)
	}

	data := map[int32]string{
		3: "ts3",
		4: "ts4",
	}

	//Очищаем хранилище
	err = memcache.FlushAll()
	if err != nil {
		t.Error(err)
	}

	for key, val := range data {
		err := memcache.Set(strconv.Itoa(int(key)), val, 0)
		if err != nil {
			t.Error(err)
		}
	}

	tests := []struct {
		name    string
		key     int32
		value   string
		wantErr bool
	}{
		{
			"1",
			1,
			"ts1",
			false,
		},
		{
			"2",
			2,
			"ts2",
			false,
		},
		{
			"3",
			3,
			"ts3",
			true,
		},
		{
			"4",
			4,
			"test",
			true,
		},
		{
			"5",
			5,
			"ts5",
			false,
		},
		{
			"4",
			100,
			"ts100",
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := MemcachedCurrentStorage.Put(context.Background(), tt.key, tt.value); (err != nil) != tt.wantErr {
				t.Errorf("Put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}

	keys := []int32{1, 2, 3, 4, 5, 100}
	for _, key := range keys {
		val, err := memcache.Get(strconv.Itoa(int(key)))
		if err != nil {
			t.Error(err)
		}

		if val != "ts"+strconv.Itoa(int(key)) {
			t.Error("value is wrong")
		}
	}
}

func TestMemConcurrency(t *testing.T) {
	err := CreateMemcahed(context.Background(), memcahedAdr)
	if err != nil {
		t.Error(err)
	}

	data := map[int32]string{
		1: "ts1",
		2: "ts2",
		3: "ts3",
	}

	//Очищаем хранилище
	err = memcache.FlushAll()
	if err != nil {
		t.Error(err)
	}

	for key, val := range data {
		err := memcache.Set(strconv.Itoa(int(key)), val, 0)
		if err != nil {
			t.Error(err)
		}
	}

	var wg sync.WaitGroup
	for i := 0; i < 1003; i++ {
		wg.Add(1)
		num := i
		go func() {
			MemcachedCurrentStorage.Put(context.Background(), int32(num), "ts"+strconv.Itoa(num))
			wg.Done()
		}()
	}

	for i := 4; i < 1000; i++ {
		wg.Add(1)
		num := i
		go func() {
			MemcachedCurrentStorage.Delete(context.Background(), int32(num))
			wg.Done()
		}()
	}

	wg.Wait()

	keys := []int32{1, 2, 3, 1000, 1001, 1002}
	for _, key := range keys {
		val, err := memcache.Get(strconv.Itoa(int(key)))
		if err != nil {
			t.Error(err)
		}

		if val != "ts"+strconv.Itoa(int(key)) {
			t.Error("value is wrong")
		}
	}

}
