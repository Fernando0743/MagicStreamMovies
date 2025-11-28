import Movie from "../movie/Movie";

const Movies = ({movies, updateMovieReview, message}) => {
  
    return (
        <div className="container mt-4">
            <div className="row">
                {/* Render movies when movies array prop is not empty and contains at least one element */}
                {movies && movies.length > 0
                    ? movies.map((movie) => (
                        // Use mongodb id for movie as key (or unique identifier)
                        <Movie key={movie.id} updateMovieReview={updateMovieReview} movie={movie} />
                    ))
                    : <h2>{message}</h2>
                }

            </div>
        </div>
    )

}

export default Movies;