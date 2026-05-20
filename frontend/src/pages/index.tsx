import Head from 'next/head';
import Link from 'next/link';

export default function Home() {
  return (
    <>
      <Head>
        <title>RecallFlow - Never Lose a Patient Because of a Missed Call</title>
        <meta name="description" content="AI-powered missed call recovery for healthcare clinics" />
      </Head>

      <div className="min-h-screen bg-gradient-to-b from-white to-gray-50">
        {/* Navigation */}
        <nav className="bg-white shadow-sm">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="flex justify-between h-16 items-center">
              <div className="flex-shrink-0 flex items-center">
                <h1 className="text-2xl font-bold text-primary-600">RecallFlow</h1>
              </div>
              <div className="flex items-center space-x-4">
                <Link href="/login" className="text-gray-700 hover:text-primary-600">
                  Login
                </Link>
                <Link href="/register" className="btn-primary">
                  Get Started
                </Link>
              </div>
            </div>
          </div>
        </nav>

        {/* Hero Section */}
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 py-20">
          <div className="text-center">
            <h2 className="text-5xl font-extrabold text-gray-900 mb-6">
              Never Lose a Patient<br />
              Because of a Missed Call
            </h2>
            <p className="text-xl text-gray-600 mb-8 max-w-3xl mx-auto">
              RecallFlow automatically recovers missed patient calls using SMS automation 
              and AI-powered conversations, helping your clinic capture every opportunity.
            </p>
            <div className="flex justify-center space-x-4">
              <Link href="/register" className="btn-primary text-lg px-8 py-3">
                Start Free Trial
              </Link>
              <a href="#features" className="btn-secondary text-lg px-8 py-3">
                Learn More
              </a>
            </div>
          </div>
        </div>

        {/* Stats Section */}
        <div className="bg-primary-600 text-white py-16">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <div className="grid grid-cols-1 md:grid-cols-3 gap-8 text-center">
              <div>
                <div className="text-4xl font-bold mb-2">85%</div>
                <div className="text-primary-100">Response Rate</div>
              </div>
              <div>
                <div className="text-4xl font-bold mb-2">$200+</div>
                <div className="text-primary-100">Avg. Recovered Per Lead</div>
              </div>
              <div>
                <div className="text-4xl font-bold mb-2">24/7</div>
                <div className="text-primary-100">Automatic Coverage</div>
              </div>
            </div>
          </div>
        </div>

        {/* Features Section */}
        <div id="features" className="py-20">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <h3 className="text-3xl font-bold text-center mb-12">How It Works</h3>
            <div className="grid md:grid-cols-3 gap-8">
              <div className="card text-center">
                <div className="text-primary-600 text-4xl mb-4">📞</div>
                <h4 className="text-xl font-semibold mb-3">Missed Call Detected</h4>
                <p className="text-gray-600">
                  When a patient call goes unanswered, RecallFlow automatically triggers 
                  the recovery workflow within seconds.
                </p>
              </div>
              <div className="card text-center">
                <div className="text-primary-600 text-4xl mb-4">💬</div>
                <h4 className="text-xl font-semibold mb-3">Automatic SMS Sent</h4>
                <p className="text-gray-600">
                  An immediate SMS is sent to the patient, keeping them engaged and 
                  showing your clinic cares.
                </p>
              </div>
              <div className="card text-center">
                <div className="text-primary-600 text-4xl mb-4">🤖</div>
                <h4 className="text-xl font-semibold mb-3">AI Handles Conversation</h4>
                <p className="text-gray-600">
                  AI classifies patient intent and provides smart responses, notifying 
                  staff when needed.
                </p>
              </div>
            </div>
          </div>
        </div>

        {/* Pricing Section */}
        <div className="bg-gray-100 py-20">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
            <h3 className="text-3xl font-bold text-center mb-12">Simple, Transparent Pricing</h3>
            <div className="grid md:grid-cols-3 gap-8">
              <div className="card">
                <h4 className="text-2xl font-bold mb-4">Solo Clinic</h4>
                <div className="text-4xl font-bold text-primary-600 mb-4">$99<span className="text-lg text-gray-600">/mo</span></div>
                <ul className="space-y-3 mb-6">
                  <li className="flex items-center">
                    <span className="text-green-500 mr-2">✓</span>
                    Unlimited missed calls
                  </li>
                  <li className="flex items-center">
                    <span className="text-green-500 mr-2">✓</span>
                    Automatic SMS replies
                  </li>
                  <li className="flex items-center">
                    <span className="text-green-500 mr-2">✓</span>
                    AI intent classification
                  </li>
                  <li className="flex items-center">
                    <span className="text-green-500 mr-2">✓</span>
                    1 location
                  </li>
                </ul>
                <Link href="/register" className="btn-primary w-full text-center block">
                  Get Started
                </Link>
              </div>
              
              <div className="card border-2 border-primary-600 relative">
                <div className="absolute top-0 right-0 bg-primary-600 text-white px-3 py-1 text-sm rounded-bl">
                  Popular
                </div>
                <h4 className="text-2xl font-bold mb-4">Multi-Provider</h4>
                <div className="text-4xl font-bold text-primary-600 mb-4">$299<span className="text-lg text-gray-600">/mo</span></div>
                <ul className="space-y-3 mb-6">
                  <li className="flex items-center">
                    <span className="text-green-500 mr-2">✓</span>
                    Everything in Solo
                  </li>
                  <li className="flex items-center">
                    <span className="text-green-500 mr-2">✓</span>
                    Up to 3 providers
                  </li>
                  <li className="flex items-center">
                    <span className="text-green-500 mr-2">✓</span>
                    Advanced analytics
                  </li>
                  <li className="flex items-center">
                    <span className="text-green-500 mr-2">✓</span>
                    Priority support
                  </li>
                </ul>
                <Link href="/register" className="btn-primary w-full text-center block">
                  Get Started
                </Link>
              </div>
              
              <div className="card">
                <h4 className="text-2xl font-bold mb-4">Multi-Location</h4>
                <div className="text-4xl font-bold text-primary-600 mb-4">$999<span className="text-lg text-gray-600">/mo</span></div>
                <ul className="space-y-3 mb-6">
                  <li className="flex items-center">
                    <span className="text-green-500 mr-2">✓</span>
                    Everything in Multi-Provider
                  </li>
                  <li className="flex items-center">
                    <span className="text-green-500 mr-2">✓</span>
                    Unlimited locations
                  </li>
                  <li className="flex items-center">
                    <span className="text-green-500 mr-2">✓</span>
                    Custom integrations
                  </li>
                  <li className="flex items-center">
                    <span className="text-green-500 mr-2">✓</span>
                    Dedicated support
                  </li>
                </ul>
                <Link href="/register" className="btn-primary w-full text-center block">
                  Get Started
                </Link>
              </div>
            </div>
          </div>
        </div>

        {/* CTA Section */}
        <div className="py-20">
          <div className="max-w-4xl mx-auto text-center px-4">
            <h3 className="text-3xl font-bold mb-6">Ready to Stop Losing Patients?</h3>
            <p className="text-xl text-gray-600 mb-8">
              Start your 14-day free trial today. No credit card required.
            </p>
            <Link href="/register" className="btn-primary text-lg px-8 py-3">
              Start Free Trial
            </Link>
          </div>
        </div>

        {/* Footer */}
        <footer className="bg-gray-900 text-white py-12">
          <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 text-center">
            <p className="text-gray-400">
              © 2026 RecallFlow Technologies LLC. All rights reserved.
            </p>
          </div>
        </footer>
      </div>
    </>
  );
}
