import Button from 'react-bootstrap/Button'
import Container from 'react-bootstrap/Container'
import Nav from 'react-bootstrap/Nav'
import Navbar from 'react-bootstrap/Navbar'
import {useNavigate, NavLink, Link} from 'react-router-dom'
import { useState } from 'react'
import useAuth from '../../hook/useAuth';
import logo from '../../assets/MagicStreamLogo.png'

const Header = ({handleLogout}) => {
    const navigate =  useNavigate();
    const {auth} = useAuth();

    return (
        <Navbar bg='dark' variant='dark' expand="lg" sticky='top' className='shadow-sm'>
            <Container>
                <Navbar.Brand>
                    <img
                        alt=''
                        src={logo}
                        width="30"
                        height="30"
                        className='d-inline-block align-top me-2'
                    />
                </Navbar.Brand>
           

            {/* Toggle helps to minimize nav bar, useful for small screens */}
            <Navbar.Toggle aria-controls='main-navbar-nav'/>

            {/* Links */}
                <Navbar.Collapse>
                    <Nav className='me-auto'>
                        <Nav.Link as = {NavLink} to={"/"}>
                            Home
                        </Nav.Link>

                        <Nav.Link as = {NavLink} to={"/recommended"}>
                            Recommended
                        </Nav.Link>
                    </Nav>

                
                
                    {/* If user is logged in we display logout button, if not we display login and register buttons */}
                    <Nav className='ms-auto align-items-center'>
                        {auth ? (
                        <>
                            <span className='me-3 text-light'>
                                Hello, <strong>{auth.first_name}</strong>
                            </span>

                            <Button variant='outline-light' size='sm' onClick={handleLogout}>
                                Logout
                            </Button>
                        </>
                        ):(
                            <>
                                <Button
                                    variant="outline-info"
                                    size="sm"
                                    className="me-2"
                                    onClick={() => navigate("/login")}
                                >
                                    Login
                                </Button>

                                <Button
                                    variant="outline-info"
                                    size="sm"
                                    className="me-2"
                                    onClick={() => navigate("/register")}
                                >
                                    Register
                                </Button>
                            </>
                        )}
                    </Nav>
                </Navbar.Collapse>
            </Container>
        </Navbar>

    )
}

export default Header;