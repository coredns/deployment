class Coredns < Formula
  desc "DNS server that chains plugins"
  homepage "https://coredns.io"
  url "https://github.com/coredns/coredns/releases/download/v1.6.9/coredns_1.6.9_darwin_amd64.tgz"
  version "1.6.9"
  sha256 "d0b91c2423e459b6c03561d640c2c686b3168f45a9c510782c268dd549f4a84f"
  head "https://github.com/coredns/coredns.git"

  def default_coredns_config; <<~EOS
    . {
      hosts {
        fallthrough
      }
      forward . https://8.8.8.8:53 https://8.8.4.4:53
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
    <string>#{opt_bin}/coredns</string>
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
    assert_match "CoreDNS-#{version}", shell_output("#{bin}/coredns -version")
  end
end
