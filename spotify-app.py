import streamlit as st
import requests
import json

st.title('Spotify Music Player')

# Function to get a random song from your server
def get_random_song():
    url = 'http://localhost:8080/spotify/random-song'  # Adjust this as necessary
    response = requests.get(url)
    if response.status_code == 200:
        return response.text, None  # Assume response.text is the song name
    else:
        return None, "Failed to fetch song"
    
def get_preview_url():
    url = 'http://localhost:8080/spotify/getPreviewURL'
    response = requests.get(url)
    print(f"HTTP status code: {response.status_code}")  # Print the HTTP status code
    if response.status_code == 200:
        data = json.loads(response.text)
        return data['preview_url'], None
    else:
        return None, "Failed to fetch preview URL"

# Button to fetch and display the song
if st.button('Generate Random Song'):
    song_name, error = get_random_song()
    preview_url, error_preview = get_preview_url()
    if error:
        st.error(error)
    elif error_preview:
        st.error(error_preview)
    else:
        if song_name:
            st.write(f"{song_name}")  # Display the song name
            if preview_url:
                st.audio(preview_url)
            else:
                st.error("No preview URL was received.")
        else:
            st.error("No song name was received.")
