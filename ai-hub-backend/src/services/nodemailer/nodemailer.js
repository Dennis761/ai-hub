import nodemailer from 'nodemailer';
import loadEnv from '../../config/loadEnv.js';

async function sendVerificationEmail(email, code) {
  try {
    const transporter = nodemailer.createTransport({
      service: 'Gmail',
      auth: {
        user: loadEnv.EMAIL_USER,
        pass: loadEnv.EMAIL_PASS,
      },
    });

    const mailOptions = {
      from: '"AI Hub" <noreply@aihub.com>',
      to: email,
      subject: 'Registration Confirmation',
      text: `Your verification code is: ${code}`,
    };

    await transporter.verify();
    console.log('SMTP connection successful');

    await transporter.sendMail(mailOptions);
    console.log('Email sent successfully');
  } catch (error) {
    console.error(' Error sending email:', error);
    throw new Error('Error while sending email: ' + error.message);
  }
}

export default sendVerificationEmail;
