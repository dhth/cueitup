package types

import (
	"strings"
	"testing"

	sqstypes "github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetJSONMessage(t *testing.T) {
	messageID := "7bc4bd4a-f099-4831-952d-5d03006a6a6f"
	// the keys need to be sorted alphabetically
	messageBody := strings.TrimSpace(`
{
  "browserInfo": {
    "browserName": "Firefox",
    "browserVersion": 118,
    "deviceType": "Desktop",
    "platform": "Linux"
  },
  "isBot": true,
  "metadata": "{\"aggregateId\":\"00000000-0000-0000-0000-000000012363\",\"sequenceNr\":347}",
  "sessionId": "987e6543-b21a-34c5-d678-123456789abc",
  "transactionId": "123e4567-e89b-12d3-a456-426614174000"
}
`)
	keyBrowserInfo := "browserInfo"
	valueBrowserInfoJSON := strings.TrimSpace(`
{
  "browserName": "Firefox",
  "browserVersion": 118,
  "deviceType": "Desktop",
  "platform": "Linux"
}
`)
	keyMetadata := "metadata"
	valueMetadataJSON := strings.TrimSpace(`
{
  "aggregateId": "00000000-0000-0000-0000-000000012363",
  "sequenceNr": 347
}
`)
	keyAggregateID := "aggregateId"
	valueAggregateID := "00000000-0000-0000-0000-000000012363"
	keySequenceNr := "sequenceNr"
	keySessionID := "sessionId"
	valueSessionID := "987e6543-b21a-34c5-d678-123456789abc"
	keyPlatform := "platform"
	valuePlatform := "Linux"
	absentKey := "absent"
	keyIsBot := "isBot"
	keyBrowserVersion := "browserVersion"
	invalidJSONBody := `not valid json`

	testCases := []struct {
		name       string
		message    sqstypes.Message
		subsetKey  *string
		contextKey *string
		expected   Message
	}{
		{
			name: "empty message ID",
			message: sqstypes.Message{
				Body: &messageBody,
			},
			expected: Message{
				Err: errMessageIDNil,
			},
		},
		{
			name: "empty message body",
			message: sqstypes.Message{
				MessageId: &messageID,
			},
			expected: Message{
				Err: errMessageBodyNil,
			},
		},
		{
			name: "invalid json body",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &invalidJSONBody,
			},
			expected: Message{
				Err: errCouldntUnmarshalBytes,
			},
		},
		{
			name: "correct json body, no subset, no context",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			expected: Message{
				ID:   messageID,
				Body: messageBody,
			},
		},
		{
			name: "correct json body, no subset, correct context",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			contextKey: &keySessionID,
			expected: Message{
				ID:           messageID,
				Body:         messageBody,
				ContextKey:   &keySessionID,
				ContextValue: &valueSessionID,
			},
		},
		{
			name: "correct json body, no subset, absent context",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			contextKey: &absentKey,
			expected: Message{
				Err: errContextKeyNotFound,
			},
		},
		{
			name: "correct json body, no subset, incorrect context key that points to an object",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			contextKey: &keyBrowserInfo,
			expected: Message{
				Err: errContextValueTypeUnsupported,
			},
		},
		{
			name: "correct json body, no subset, incorrect context key that points to a bool",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			contextKey: &keyIsBot,
			expected: Message{
				Err: errContextValueTypeUnsupported,
			},
		},
		{
			name: "correct json body, correct subset, no context key",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			subsetKey: &keyBrowserInfo,
			expected: Message{
				ID:   messageID,
				Body: valueBrowserInfoJSON,
			},
		},
		{
			name: "correct json body, absent subset, no context key",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			subsetKey: &absentKey,
			expected: Message{
				Err: errSubsetKeyNotFound,
			},
		},
		{
			name: "correct json body, incorrect subset key that points to a bool, no context key",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			subsetKey: &keyIsBot,
			expected: Message{
				Err: errSubsetTypeIsUnsupported,
			},
		},
		{
			name: "correct json body, incorrect subset key that points to a non stringified JSON, no context key",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			subsetKey: &keySessionID,
			expected: Message{
				Err: errCouldntUnmarshalSubsetValue,
			},
		},
		{
			name: "correct json body, correct subset, correct context key",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			subsetKey:  &keyBrowserInfo,
			contextKey: &keyPlatform,
			expected: Message{
				ID:           messageID,
				Body:         valueBrowserInfoJSON,
				ContextKey:   &keyPlatform,
				ContextValue: &valuePlatform,
			},
		},
		{
			name: "correct json body, correct subset, absent context key",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			subsetKey:  &keyBrowserInfo,
			contextKey: &absentKey,
			expected: Message{
				Err: errContextKeyNotFound,
			},
		},
		{
			name: "correct json body, correct subset, incorrect context key that points to a float",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			subsetKey:  &keyBrowserInfo,
			contextKey: &keyBrowserVersion,
			expected: Message{
				Err: errContextValueTypeUnsupported,
			},
		},
		{
			name: "correct json body, correct subset key that points to stringified JSON, no context key",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			subsetKey: &keyMetadata,
			expected: Message{
				ID:   messageID,
				Body: valueMetadataJSON,
			},
		},
		{
			name: "correct json body, correct subset key that points to stringified JSON, correct context key",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			subsetKey:  &keyMetadata,
			contextKey: &keyAggregateID,
			expected: Message{
				ID:           messageID,
				Body:         valueMetadataJSON,
				ContextKey:   &keyAggregateID,
				ContextValue: &valueAggregateID,
			},
		},
		{
			name: "correct json body, correct subset key that points to stringified JSON, absent context key",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			subsetKey:  &keyMetadata,
			contextKey: &absentKey,
			expected: Message{
				Err: errContextKeyNotFound,
			},
		},
		{
			name: "correct json body, correct subset key that points to stringified JSON, incorrect context key that points to a float",
			message: sqstypes.Message{
				MessageId: &messageID,
				Body:      &messageBody,
			},
			subsetKey:  &keyMetadata,
			contextKey: &keySequenceNr,
			expected: Message{
				Err: errContextValueTypeUnsupported,
			},
		},
	}

	for _, tt := range testCases {
		got := getJSONMessage(&tt.message, tt.subsetKey, tt.contextKey)

		if tt.expected.Err == nil {
			assert.Equal(t, tt.expected, got, tt.name)
		} else {
			require.ErrorIs(t, got.Err, tt.expected.Err, tt.name)
			assert.Equal(t, got.ID, tt.expected.ID, tt.name)
			assert.Equal(t, got.Body, tt.expected.Body, tt.name)
			assert.Equal(t, got.ContextKey, tt.expected.ContextKey, tt.name)
			assert.Equal(t, got.ContextValue, tt.expected.ContextValue, tt.name)
		}
	}
}
