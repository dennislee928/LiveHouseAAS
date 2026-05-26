package payment

import (
	"context"
	"testing"
)

func TestMockProvider(t *testing.T) {
	p := NewMockProvider()

	if p.Name() != ProviderMock {
		t.Errorf("expected mock, got %s", p.Name())
	}

	req := &PaymentRequest{
		OrderID:     "test-order-1",
		Amount:      1000,
		Currency:    "TWD",
		Description: "Test payment",
	}

	result, err := p.Pay(context.Background(), req)
	if err != nil {
		t.Fatalf("Pay failed: %v", err)
	}

	if result.Status != "completed" {
		t.Errorf("expected completed, got %s", result.Status)
	}
}

func TestMockRefund(t *testing.T) {
	p := NewMockProvider()
	err := p.Refund(context.Background(), "tx-123", 500)
	if err != nil {
		t.Errorf("Refund failed: %v", err)
	}
}

func TestRouter(t *testing.T) {
	r := NewRouter()
	if r == nil {
		t.Fatal("router is nil")
	}

	_, err := r.Pay(context.Background(), ProviderMock, &PaymentRequest{
		OrderID:  "test",
		Amount:   500,
		Currency: "TWD",
	})
	if err != nil {
		t.Fatalf("Pay via router failed: %v", err)
	}

	// Unregistered provider should fail
	_, err = r.Pay(context.Background(), "unknown", &PaymentRequest{})
	if err == nil {
		t.Error("expected error for unknown provider")
	}
}

func TestECPayCheckMacValue(t *testing.T) {
	p := NewECPayProvider("2000132", "hashKey", "hashIV", true)

	form := make(map[string]string)
	form["MerchantID"] = "2000132"
	form["MerchantTradeNo"] = "TEST123"
	form["TotalAmount"] = "1000"

	// Just verify it doesn't panic
	_ = p.VerifyCheckMacValue(form)
}

func TestNewebPayEncryptDecrypt(t *testing.T) {
	p := NewNewebPayProvider("MS12345678", "hashKey12345678", "hashIV12345678", true)

	plaintext := `{"MerchantID":"MS12345678","Amt":100}`
	encrypted := p.encryptAES(plaintext)
	if encrypted == "" {
		t.Error("encryption produced empty string")
	}

	decrypted, err := p.decryptAES(encrypted)
	if err != nil {
		t.Fatalf("decryption failed: %v", err)
	}

	if decrypted != plaintext {
		t.Errorf("decrypted mismatch: got %s, want %s", decrypted, plaintext)
	}
}

func TestPKCS7Pad(t *testing.T) {
	data := []byte{1, 2, 3}
	padded := pkcs7Pad(data, 16)
	if len(padded) != 16 {
		t.Errorf("expected 16 bytes, got %d", len(padded))
	}
}

func TestNewebPaySHA256(t *testing.T) {
	p := NewNewebPayProvider("MS12345678", "key", "iv", true)
	hash := p.sha256Hex("test")
	if hash == "" {
		t.Error("sha256 hex produced empty string")
	}
}
