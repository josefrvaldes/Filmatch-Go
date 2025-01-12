import requests
import random

# API URLs
BASE_API_URL = "http://localhost:8080"
USER_API_URL = f"{BASE_API_URL}/user"  # Endpoint to create users
MOVIE_API_URL = f"{BASE_API_URL}/user/content"  # Endpoint to post user content
API_KEY = "your-api-key"  

# Create a user in the API
def create_user(user):
    headers = {
        "Content-Type": "application/json",
        "X-API-KEY": API_KEY  
    }
    try:
        response = requests.post(USER_API_URL, json=user, headers=headers)
        if response.status_code == 201:  # HTTP 201 Created
            print(f"User created: {user['email']}")
            return True
        elif response.status_code == 409:  # HTTP 409 Conflict (existing user)
            print(f"User already exists: {user['email']}")
            return True
        else:
            print(f"Failed to create user ({response.status_code}): {response.text}")
            return False
    except Exception as e:
        print(f"Error creating user: {e}")
        return False

# Generate test users
def generate_user(index):
    return {
        "email": f"user{index}@example.com",
        "username": f"user{index}"
    }

# Generate test movies
def generate_movie(index):
    return {
        "id": 1000 + index,
        "title": f"Movie {index}",
        "original_title": f"Original Movie {index}",
        "overview": f"This is a test movie description for Movie {index}.",
        "release_date": f"202{random.randint(0, 3)}-0{random.randint(1, 9)}-1{random.randint(0, 9)}",
        "popularity": round(random.uniform(1.0, 100.0), 2),
        "vote_average": round(random.uniform(1.0, 10.0), 1),
        "vote_count": random.randint(10, 1000),
        "adult": False,
        "backdrop_path": f"/backdrop_movie_{index}.jpg",
        "poster_path": f"/poster_movie_{index}.jpg",
        "genre_ids": [random.randint(1, 20), random.randint(21, 40)]
    }

# Generate test TV shows
def generate_tv_show(index):
    return {
        "id": 2000 + index,
        "name": f"TV Show {index}",
        "original_name": f"Original TV Show {index}",
        "overview": f"This is a test TV show description for TV Show {index}.",
        "first_air_date": f"202{random.randint(0, 3)}-0{random.randint(1, 9)}-1{random.randint(0, 9)}",
        "popularity": round(random.uniform(1.0, 100.0), 2),
        "vote_average": round(random.uniform(1.0, 10.0), 1),
        "vote_count": random.randint(10, 1000),
        "adult": False,
        "backdrop_path": f"/backdrop_tv_show_{index}.jpg",
        "poster_path": f"/poster_tv_show_{index}.jpg",
        "origin_country": [random.choice(["US", "UK", "KR", "JP"])]
    }

# Let's post some data to the API
def post_to_api(user, content, is_movie=True, status=random.randint(1, 4)):
    payload = {
        "user": user,
        "status": status
    }
    if is_movie:
        payload["movie"] = content
    else:
        payload["tv_show"] = content

    headers = {
        "Content-Type": "application/json",
        "X-API-KEY": API_KEY  # Include API key if your API requires it
    }

    try:
        response = requests.post(MOVIE_API_URL, json=payload, headers=headers)
        if response.status_code == 200:
            print(f"Success: {payload['user']['email']} -> {content['id']}")
        else:
            print(f"Failed ({response.status_code}): {response.text}")
    except Exception as e:
        print(f"Error: {e}")

# Populate the API with some data
def populate_api():
    num_users = 10
    num_movies = 5
    num_tv_shows = 5

    for user_index in range(1, num_users + 1):
        user = generate_user(user_index)

        # Create user
        if not create_user(user):
            print(f"Skipping user {user['email']} due to creation failure")
            continue

        # Create movies for this user
        for movie_index in range(1, num_movies + 1):
            movie = generate_movie(movie_index)
            post_to_api(user, movie, is_movie=True)

        # Create TV shows for this user
        for tv_index in range(1, num_tv_shows + 1):
            tv_show = generate_tv_show(tv_index)
            post_to_api(user, tv_show, is_movie=False)

if __name__ == "__main__":
    populate_api()