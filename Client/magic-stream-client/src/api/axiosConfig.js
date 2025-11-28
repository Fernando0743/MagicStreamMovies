import axios from 'axios'

//Use axios package to read serverside code (or our backend code)
//Base URL of our backend (server-side) code
const apiURL = import.meta.env.VITE_API_BASE_URL;

export default axios.create({
    baseURL: apiURL,
    headers: {'Content-Type':'application/json'},
    //http cookies
    withCredentials: true
})