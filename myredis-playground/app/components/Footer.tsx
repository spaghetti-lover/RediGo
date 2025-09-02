export default function Footer() {
  return (
    <footer className="w-full bg-gruv-dark/50 backdrop-blur-sm">
      <div className="max-w-6xl mx-auto px-6 py-4">
        <div className="flex flex-col sm:flex-row items-center justify-between gap-4">
          {/* License Info */}
          <div className="flex items-center gap-2 text-sm text-gruv-muted">
            <span>©</span>
            <span>MIT License</span>
            <span className="text-gruv-gray">•</span>
            <span>Open Source</span>
          </div>

          {/* Star Call-to-Action */}
          <div className="flex items-center gap-3">
            <span className="text-sm text-gruv-muted">Like this project?</span>
            <a
              href="https://github.com/spaghetti-lover/RediGo"
              target="_blank"
              rel="noopener noreferrer"
              className="flex items-center gap-2 px-4 py-2 bg-gruv-yellow text-gruv-dark rounded-lg hover:bg-gruv-yellow/90 transition-all duration-200 font-medium text-sm group"
            >
              <svg
                className="w-4 h-4 group-hover:scale-110 transition-transform duration-200"
                fill="currentColor"
                viewBox="0 0 20 20"
              >
                <path d="M9.049 2.927c.3-.921 1.603-.921 1.902 0l1.07 3.292a1 1 0 00.95.69h3.462c.969 0 1.371 1.24.588 1.81l-2.8 2.034a1 1 0 00-.364 1.118l1.07 3.292c.3.921-.755 1.688-1.54 1.118l-2.8-2.034a1 1 0 00-1.175 0l-2.8 2.034c-.784.57-1.838-.197-1.539-1.118l1.07-3.292a1 1 0 00-.364-1.118L2.98 8.72c-.783-.57-.38-1.81.588-1.81h3.461a1 1 0 00.951-.69l1.07-3.292z" />
              </svg>
              <span className="hover:text-yellow-400">Star on GitHub</span>
            </a>
          </div>
        </div>
      </div>
    </footer>
  );
}
