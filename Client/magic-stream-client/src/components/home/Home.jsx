import { useState, useEffect } from "react";
import axiosClient from '..//..//api/axiosConfig'
import Movies from '../movies/Movies'

const Home = ({updateMovieReview}) => {
    const [movies, setMovies] = useState([]);
    const [loading, setLoading] = useState(false)
    const [message, setMessage] = useState()

    //When Home component loads, we use UseEffect to call the movies endpoint on the server and populate movies collection with 
    //movie data from the server
    useEffect(() => {
        const fetchMovies = async () => {
            //Display loading indicator while servserside code is retrieving relevant data
            setLoading(true);
            setMessage("");
            try{
                //Call serverside endpoint
                const response = await axiosClient.get('/movies');
                setMovies(response.data);
                if(response.data.length == 0){
                    setMessage("There are currently no movies available")
                }
            }catch(error){
                console.error('Error fetching movies')
                setMessage("Error fetching movies")
            }finally{
                setLoading(false)
            }
        }
        fetchMovies();
    }, [])

    return (
        <>
            {loading ? (
                <h2>Loading...</h2>
            ) : (
                <Movies movies= {movies} updateMovieReview = {updateMovieReview}  message={message}/>
            )}
        </>
    );

    
};

export default Home