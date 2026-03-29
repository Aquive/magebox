# typed: false
# frozen_string_literal: true

class Magebox < Formula
  desc "Fast, native Magento development environment"
  homepage "https://magebox.dev"
  version "1.2.7"
  license "MIT"

  on_macos do
    on_arm do
      url "https://github.com/qoliber/magebox/releases/download/v#{version}/magebox-darwin-arm64"
      sha256 "PLACEHOLDER_DARWIN_ARM64"
    end
    on_intel do
      url "https://github.com/qoliber/magebox/releases/download/v#{version}/magebox-darwin-amd64"
      sha256 "PLACEHOLDER_DARWIN_AMD64"
    end
  end

  on_linux do
    on_arm do
      url "https://github.com/qoliber/magebox/releases/download/v#{version}/magebox-linux-arm64"
      sha256 "PLACEHOLDER_LINUX_ARM64"
    end
    on_intel do
      url "https://github.com/qoliber/magebox/releases/download/v#{version}/magebox-linux-amd64"
      sha256 "PLACEHOLDER_LINUX_AMD64"
    end
  end

  def install
    suffix = if OS.mac?
               Hardware::CPU.arm? ? "darwin-arm64" : "darwin-amd64"
             else
               Hardware::CPU.arm? ? "linux-arm64" : "linux-amd64"
             end
    bin.install "magebox-#{suffix}" => "magebox"
  end

  test do
    assert_match "MageBox", shell_output("#{bin}/magebox --version")
  end
end
