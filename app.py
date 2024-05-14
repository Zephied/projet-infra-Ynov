import streamlit as st
import requests


st.title("ChatGPT-like clone")

if 'messages' not in st.session_state:
    st.session_state['messages'] = [
        {"role": "assistant", "content": "You are a helpful assistant."}
    ]


st.sidebar.title("Sidebar")
clear_button = st.sidebar.button("Clear Conversation", key="clear")
# reset everything
if clear_button:
    st.session_state['messages'] = [
        {"role": "system", "content": "You are a helpful assistant."}
    ]
    
st.sidebar.subheader("About App")
st.sidebar.text("This is a ChatGPT clone app created \nusing GOLANG AND PYTHON STREAMLIT \n(FOR FRONTEND OF THE APPLICATION)")

if "messages" not in st.session_state:
    st.session_state.messages = []

for message in st.session_state.messages:
    with st.chat_message(message["role"]):
        st.markdown(message["content"])

if prompt := st.chat_input("What is up?"):
    st.session_state.messages.append({"role": "user", "content": prompt})
    with st.chat_message("user"):
        st.markdown(prompt)

    with st.chat_message("assistant"):
        response = requests.post('http://localhost:8080/chat', json={"prompt": prompt})
        if response.status_code == 200:
            data = response.json()   
        st.write(data["answer"])
    st.session_state.messages.append({"role": "assistant", "content": data["answer"]})
    
    
