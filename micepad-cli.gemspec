require_relative "lib/micepad/version"

Gem::Specification.new do |spec|
  spec.name          = "micepad-cli"
  spec.version       = Micepad::VERSION
  spec.authors       = ["Micepad"]
  spec.email         = ["support@micepad.co"]

  spec.summary       = "CLI for Micepad event management platform"
  spec.description   = "Command-line interface for managing events, participants, check-ins, and more on Micepad."
  spec.homepage      = "https://micepad.co"
  spec.license       = "MIT"

  spec.required_ruby_version = ">= 3.0"

  spec.files         = Dir["lib/**/*", "bin/*", "LICENSE.txt", "README.md"]
  spec.bindir        = "bin"
  spec.executables   = ["micepad"]

  spec.add_dependency "terminalwire-client", "~> 0.3"
end
