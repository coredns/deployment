class Coredns < Formula
  desc "DNS server that chains plugins"
  homepage "https://coredns.io"
  url "https://github.com/coredns/coredns/releases/download/v1.5.2/coredns_1.5.2_darwin_amd64.tgz"
  version "1.5.2"
  sha256 "195c73f6ccbb013b7a326d46cd55c5e0aaa66914cfab97b2030c67ddf2d7b25a"
  head "https://github.com/coredns/coredns.git"

  def default_coredns_config; <<~EOS
    . {
      hosts {
        fallthrough
      }
      proxy . 8.8.8.8:53 8.8.4.4:53 {
        protocol https_google
      }
      cache
      errors
    }
    EOS
  end

  def install
    (buildpath/"Corefile.example").write default_coredns_config
    (etc/"coredns").mkpath
    etc.install "Corefile.example" => "coredns/Corefile"
    bin.install "coredns"
  end

  def caveats; <<~EOS
    To configure coredns, take the default configuration at
    #{etc}/coredns/Corefile and edit to taste.

    By default it is configured to proxy all dns requests
    through Google's DNS-over-HTTPS:
    (https://developers.google.com/speed/public-dns/docs/dns-over-https).
    EOS
  end

  plist_options :startup => true

  def plist; <<~EOS
    <?xml version="1.0" encoding="UTF-8"?>
    <!DOCTYPE plist PUBLIC "-//Apple Computer//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
    <plist version="1.0">
    <dict>
    <key>Label</key>
    <string>#{plist_name}</string>
    <key>ProgramArguments</key>
    <array>
    <string>#{opt_sbin}/coredns</string>
    <string>-conf</string>
    <string>#{etc}/coredns/Corefile</string>
    </array>
    <key>RunAtLoad</key>
    <true/>
    <key>KeepAlive</key>
    <true/>
    <key>StandardErrorPath</key>
    <string>#{var}/log/coredns.log</string>
    <key>StandardOutPath</key>
    <string>#{var}/log/coredns.log</string>
    <key>WorkingDirectory</key>
    <string>#{HOMEBREW_PREFIX}</string>
    </dict>
    </plist>
    EOS
  end

  test do
    assert_match "CoreDNS-#{version}", shell_output("#{sbin}/coredns -version")
  end
end
