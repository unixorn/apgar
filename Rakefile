#
# Make it easy to package and test Apgar
require 'fileutils'

PACKAGE_DESCRIPTION='Apgar is a health check driver tool'
PACKAGE_URL='https://github.com/unixorn/apgar'
SEMVER=File.open('VERSION') { |file| file.read }.chomp
iteration=`git rev-list HEAD --count`.chomp

# Of course OS X and linux can't play well with each other
if `uname`.chomp == 'Darwin'
  INSTALL_GROUP = 'admin'
end

if `uname`.chomp == 'Linux'
  INSTALL_GROUP = 'root'
end

task :default => [:usage]
task :help => [:usage]

task :usage do
  puts "rake deb:      Create an apgar deb file"
  puts "rake rpm:      Create an apgar rpm file"
  puts "rake test:     Integration test apgar-probe"
  puts
  puts "You must gem install fpm to build deb or rpm files."
end

desc "Install required gems"
task :bundle_install do
  sh %{ bundle install }
end

desc "Package apgar as a DEB"
task :deb => [:fakeroot, :apgar_binaries, :bundle_install] do
  sh %{ bundle exec fpm -s dir -t deb -n apgar \
    -v #{SEMVER} --iteration #{iteration} \
    --url #{PACKAGE_URL} \
    --description "#{PACKAGE_DESCRIPTION}" \
    -C .fakeroot --license "Public Domain" etc usr var }
end

task :apgar_binaries => [:apgar_probe, :apgar_server]

desc "Package apgar as a RPM"
task :rpm => [:fakeroot, :apgar_binaries, :bundle_install] do
  sh %{ bundle exec fpm -s dir -t rpm -n apgar \
    -v #{SEMVER} --iteration #{iteration} \
    --url #{PACKAGE_URL} \
    --description "${PACKAGE_DESCRIPTION}" \
    -C .fakeroot --license "MIT" etc usr var }
end

task :fakeroot => [:apgar_binaries] do
  sh %{ rm -fr .fakeroot }
  FileUtils::mkdir_p '.fakeroot/etc/apgar/healthchecks'
  FileUtils::mkdir_p '.fakeroot/usr/local/sbin'
  FileUtils::mkdir_p '.fakeroot/var/lib/apgar'
  sh %{ cp apgar-probe apgar-server .fakeroot/usr/local/sbin}
end

task :apgar_probe do
  sh %{ go build apgar-probe.go }
end

task :apgar_server do
  sh %{ go build apgar-server.go }
end

task :test_setup do
  sh %{ mkdir -p tmp }
  sh %{ bundle install --deployment }
end

desc "Cleanup after build"
task :cleanup do
  sh %{ find . -name '*.o' -exec rm '{}' ';' }
  sh %{ find . -name '*.un~' -exec rm '{}' ';' }
  sh %{ rm -fr .fakeroot apgar-probe apgar-server tmp *.deb }
  if File.directory?('./test/reports')
    sh %{ find ./test/reports -name '*.xml' -exec rm '{}' ';' }
  end
end

desc "Run test suite"
task :test => [:apgar_binaries, :test_setup, :bundle_install] do
  sh %{ bundle exec ./apgar_tests.rb }
end

desc "Format go files"
task :fmt do
  sh %{ go fmt *.go }
end

task :c => [:cleanup]
task :f => [:fmt]
task :t => [:test]
task :v => [:verbose_test]

desc "Verbose test"
task :verbose_test => [:apgar_probe, :test_setup] do
  sh %{ echo ziggy > tmp/status }
  system("./apgar-probe --debug 50 --document-root tmp --healthcheck-tree fixtures/005-suffix-passes")
  system("cat tmp/status")
end
