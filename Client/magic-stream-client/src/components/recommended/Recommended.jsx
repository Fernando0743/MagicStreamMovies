import useAxiosPrivate from "../../hook/useAxiosPrivate";
import { useEffect, useState } from "react";
import Movies from '../movies/Movies'


const Recommended = () => {
    const [movies, setMovies] = useState([]);
    const [loading, setLoading] = useState(false);
    const [message, setMessage] = useState();
    const axiosPrivate = useAxiosPrivate();

    //Call recommended movies endpoint when this component starts loading
    useEffect(() =>{
        const fetchRecommendedMovies = async () => {
            setLoading(true);
            setMessage("");

            try{
                //Use axios private so when we call a protected endpoint, the token s being sent along with the http request
                //(thanks to the interceptor as well)
                const response = await axiosPrivate.get('/recommendedmovies');
                setMovies(response.data)
            } catch(error){
                console.error("Error fetching recommended movies: ", error);
            } finally {
                setLoading(false);
            }
        }
        fetchRecommendedMovies();
    }, [])

    return (
        <>
        {loading ? (
            <h2>Loading....</h2>
        ) : (
            <Movies movies={movies} message={message} />
        )}
        </>
    )
}

export default Recommended