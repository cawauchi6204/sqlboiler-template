#!/bin/bash

# 基本設定
API_URL="http://localhost:8080"
TOKEN_FILE=".token"

# カラー出力用の設定
RED='\033[0;31m'
GREEN='\033[0;32m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# ヘルパー関数
print_header() {
    echo -e "\n${BLUE}=== $1 ===${NC}\n"
}

print_response() {
    if [ $1 -eq 0 ]; then
        echo -e "${GREEN}Success${NC}"
    else
        echo -e "${RED}Failed${NC}"
    fi
    echo "$2" | jq '.'
}

get_token() {
    if [ -f "$TOKEN_FILE" ]; then
        cat "$TOKEN_FILE"
    fi
}

save_token() {
    echo "$1" > "$TOKEN_FILE"
}

# 認証なしのエンドポイント

# ユーザー登録
register_user() {
    print_header "ユーザー登録"
    response=$(curl -s -X POST "$API_URL/register" \
        -H "Content-Type: application/json" \
        -d '{
            "username": "testuser",
            "display_name": "Test User",
            "email": "test@example.com",
            "password": "password123",
            "bio": "This is a test user"
        }')
    print_response $? "$response"
}

# ログイン
login() {
    print_header "ログイン"
    response=$(curl -s -X POST "$API_URL/login" \
        -H "Content-Type: application/json" \
        -d '{
            "email": "test@example.com",
            "password": "password123"
        }')
    print_response $? "$response"
    
    # トークンを保存
    token=$(echo "$response" | jq -r '.token')
    if [ "$token" != "null" ]; then
        save_token "$token"
        echo "Token saved successfully"
    fi
}

# 認証が必要なエンドポイント

# プロフィール取得
get_profile() {
    local user_id=${1:-1}
    print_header "プロフィール取得 (ID: $user_id)"
    token=$(get_token)
    response=$(curl -s -X GET "$API_URL/api/users/$user_id" \
        -H "Authorization: Bearer $token")
    print_response $? "$response"
}

# プロフィール更新
update_profile() {
    print_header "プロフィール更新"
    token=$(get_token)
    response=$(curl -s -X PUT "$API_URL/api/users/me" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        -d '{
            "display_name": "Updated Name",
            "bio": "Updated bio text",
            "profile_image_url": "https://example.com/image.jpg"
        }')
    print_response $? "$response"
}

# ツイート投稿
create_tweet() {
    print_header "ツイート投稿"
    token=$(get_token)
    response=$(curl -s -X POST "$API_URL/api/tweets" \
        -H "Authorization: Bearer $token" \
        -H "Content-Type: application/json" \
        -d '{
            "content": "This is a test tweet",
            "image_url": "https://example.com/tweet-image.jpg"
        }')
    print_response $? "$response"
}

# ツイート取得
get_tweet() {
    local tweet_id=${1:-1}
    print_header "ツイート取得 (ID: $tweet_id)"
    token=$(get_token)
    response=$(curl -s -X GET "$API_URL/api/tweets/$tweet_id" \
        -H "Authorization: Bearer $token")
    print_response $? "$response"
}

# タイムライン取得
get_timeline() {
    print_header "タイムライン取得"
    token=$(get_token)
    response=$(curl -s -X GET "$API_URL/api/tweets/timeline" \
        -H "Authorization: Bearer $token")
    print_response $? "$response"
}

# フォロー
follow_user() {
    local user_id=${1:-2}
    print_header "ユーザーをフォロー (ID: $user_id)"
    token=$(get_token)
    response=$(curl -s -X POST "$API_URL/api/users/$user_id/follow" \
        -H "Authorization: Bearer $token")
    print_response $? "$response"
}

# アンフォロー
unfollow_user() {
    local user_id=${1:-2}
    print_header "ユーザーをアンフォロー (ID: $user_id)"
    token=$(get_token)
    response=$(curl -s -X DELETE "$API_URL/api/users/$user_id/follow" \
        -H "Authorization: Bearer $token")
    print_response $? "$response"
}

# いいね
like_tweet() {
    local tweet_id=${1:-1}
    print_header "ツイートにいいね (ID: $tweet_id)"
    token=$(get_token)
    response=$(curl -s -X POST "$API_URL/api/tweets/$tweet_id/like" \
        -H "Authorization: Bearer $token")
    print_response $? "$response"
}

# いいね解除
unlike_tweet() {
    local tweet_id=${1:-1}
    print_header "いいねを解除 (ID: $tweet_id)"
    token=$(get_token)
    response=$(curl -s -X DELETE "$API_URL/api/tweets/$tweet_id/like" \
        -H "Authorization: Bearer $token")
    print_response $? "$response"
}

# メイン処理
case "$1" in
    "register")
        register_user
        ;;
    "login")
        login
        ;;
    "profile")
        get_profile $2
        ;;
    "update-profile")
        update_profile
        ;;
    "tweet")
        create_tweet
        ;;
    "get-tweet")
        get_tweet $2
        ;;
    "timeline")
        get_timeline
        ;;
    "follow")
        follow_user $2
        ;;
    "unfollow")
        unfollow_user $2
        ;;
    "like")
        like_tweet $2
        ;;
    "unlike")
        unlike_tweet $2
        ;;
    *)
        echo "使用方法:"
        echo "  $0 register                # 新規ユーザー登録"
        echo "  $0 login                   # ログイン"
        echo "  $0 profile [id]            # プロフィール取得"
        echo "  $0 update-profile          # プロフィール更新"
        echo "  $0 tweet                   # ツイート投稿"
        echo "  $0 get-tweet [id]          # ツイート取得"
        echo "  $0 timeline                # タイムライン取得"
        echo "  $0 follow [user_id]        # ユーザーをフォロー"
        echo "  $0 unfollow [user_id]      # ユーザーをアンフォロー"
        echo "  $0 like [tweet_id]         # ツイートにいいね"
        echo "  $0 unlike [tweet_id]       # いいねを解除"
        ;;
esac