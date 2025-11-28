import {useState} from 'react'
import Container from 'react-bootstrap/Container';
import Button from 'react-bootstrap/Button';
import Form from 'react-bootstrap/Form';
import axiosClient from '../../api/axiosConfig';
import { useNavigate, Link, useLocation } from 'react-router-dom';
import axios from 'axios';
import useAuth from '../../hook/useAuth';
import logo from '../../assets/MagicStreamLogo.png'



const Login = () => {
    
    const {setAuth} = useAuth();
    const [email, setEmail] = useState('');
    const [password, setPassword] = useState('');

    const [error, setError] = useState('');
    const [loading, setLoading] = useState(false);
    const location = useLocation();
    //const {setAuth} = useAuth();
    const navigate = useNavigate();

    const from = location.state?.from?.pathname || "/"

    const handleSubmit = async(e) => {
        e.preventDefault();
        setLoading(true);
        setError(null);

        try{
            const response = await axiosClient.post('/login', { email, password});
            console.log(response.data);
            if(response.data.error){
                setError(response.data.error);
                return;
            }
            //Setting auth globally
            setAuth(response.data);
            //localStorage.setItem('user', JSON.stringify(response.data));
            //Handle successful login (store token, redirect)
            //Direct the user to the screen that the user was trying to access before the user was authenticated
            navigate(from, {replace: true});
        } catch(err) {
            console.log(err);
            setError("Invalid email or password");
        }
    }


    return (
        <Container className="login-container d-flex align-items-center justify-content-center min-vh-100">
            <div className="login-card shadow p-4 rounded bg-white" style={{maxWidth: 400, width: '100%'}}>
                <div className="text-center mb-4">
                    <img src={logo} alt='Logo' width={60} className='mb-2'/>
                    <h2 className="fw-bold">Sign In</h2>
                    <p className="text-muted">Welcome back! Please login to your account.</p>
                </div>
                {error && <div className="alert alert-danger py-2">{error}</div>}
                <Form onSubmit={handleSubmit}>
                    <Form.Group controlId="formBasicEmail" className="mb-3">
                        <Form.Label>Email address</Form.Label>
                        <Form.Control
                            type="email"
                            placeholder="Enter email"
                            value={email}
                            onChange={(e) => setEmail(e.target.value)}
                            required
                            autoFocus
                        />
                    </Form.Group>

                    <Form.Group controlId="formBasicPassword" className="mb-3">
                        <Form.Label>Password</Form.Label>
                        <Form.Control
                            type="password"
                            placeholder="Password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                            required
                        />
                    </Form.Group>

                    <Button
                        variant="primary"
                        type="submit"
                        className="w-100 mb-2"
                        disabled={loading}
                        style={{fontWeight: 600, letterSpacing: 1}}
                    >
                        {loading ? (
                            <>
                                <span className="spinner-border spinner-border-sm me-2" role="status" aria-hidden="true"></span>
                                Logging in...
                            </>
                        ) : 'Login'}
                    </Button>
                </Form>
                <div className="text-center mt-3">
                    <span className="text-muted">Don't have an account? </span>
                    <Link to="/register" className="fw-semibold">Register here</Link>
                </div>
            </div>
        </Container>

    )
}

export default Login;