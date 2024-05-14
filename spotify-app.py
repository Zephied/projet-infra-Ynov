import streamlit as st
import requests

st.title('Spotify Music Player')

# Function to get a random song from your server
def get_random_song():
    url = 'http://localhost:8080/spotify/random-song'  # Adjust this as necessary
    response = requests.get(url)
    if response.status_code == 200:
        return response.text, None  # Assume response.text is the song name
    else:
        return None, "Failed to fetch song"

# Button to fetch and display the song
if st.button('Generate Random Song'):
    song_name, error = get_random_song()
    if error:
        st.error(error)
    else:
        if song_name:
            st.write(f"{song_name}")  # Display the song name
        else:
            st.error("No song name was received.")
