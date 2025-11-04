import React, { useState } from 'react';
import { BrowserRouter as Router, Route, Routes } from "react-router-dom";
import AdminLogin from './AdminLogin.jsx'
import AppLayout from './AppLayout.jsx'
import Projects from './Projects/Projects.jsx';
import ApiKeys from './ApiKeys.jsx'

function Nav() {
  return (
    <>
      <Router>
          <Routes>
            <Route path='/login' element={<AdminLogin/>}/>
            <Route path="/ai-hub" element={<AppLayout />}>
              <Route path="keys" element={<ApiKeys />} />
              <Route path="projects" element={<Projects />} />
            </Route>
          </Routes>
      </Router>
      </>
  );
}

export default Nav;