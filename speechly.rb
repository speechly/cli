# This file was generated by GoReleaser. DO NOT EDIT.
class Speechly < Formula
  desc ""
  homepage "https://www.speechly.com/"
  version "0.0.2"
  bottle :unneeded

  if OS.mac?
    url "https://github.com/speechly/cli/releases/download/v0.0.2/speechly_0.0.2_macOS_x86_64.tar.gz"
    sha256 "5c0b6360fdf9567c5343e873a59128cb4bec6ca09cb1bbfead5b8cea9dd59ed1"
  elsif OS.linux?
    if Hardware::CPU.intel?
      url "https://github.com/speechly/cli/releases/download/v0.0.2/speechly_0.0.2_Linux_x86_64.tar.gz"
      sha256 "73a91fba8177995352fd924d4161ed83d3134f4e9250df5fea8a8856cb0f256e"
    end
  end

  def install
    bin.install "speechly"
  end
end
