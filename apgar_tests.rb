#!/usr/bin/env ruby
#
require 'minitest/autorun'
require 'minitest/reporters'

Minitest::Reporters.use! [Minitest::Reporters::SpecReporter.new(),
                          Minitest::Reporters::JUnitReporter.new]

STATUS_FILE='tmp/status'

class TestApgarProbe < MiniTest::Test

  # We want to always run the tests in alphanumeric order so that the console
  # order is predictable. These tests re-run the binary every time, so there's
  # no gain from running them in a random order.
  def self.test_order
   :alpha
  end

  def setup
    # Make sure the file content checks in our tests aren't spoofed by stale status file
    File.delete(STATUS_FILE) if File.exist?(STATUS_FILE)
  end

  def test_multiple_failing
    _ = `./apgar-probe --document-root tmp --healthcheck-tree fixtures/004-multiple-failing`
    exitcode = $?.to_i
    assert_equal false, (exitcode == 0)
    assert_equal "UNHEALTHY\n", File.open(STATUS_FILE) { |file| file.read }
  end

  def test_multiple_passing
    run = `./apgar-probe --document-root tmp --healthcheck-tree fixtures/003-multiple-passing`
    exitcode = $?.to_i
    assert_equal true, (exitcode == 0)
    assert_equal "OK\n", File.open(STATUS_FILE) { |file| file.read }
  end

  def test_single_failing
    run = `./apgar-probe --document-root tmp --healthcheck-tree fixtures/002-single-failing`
    exitcode = $?.to_i
    assert_equal false, (exitcode == 0)
    assert_equal "UNHEALTHY\n", File.open(STATUS_FILE) { |file| file.read }
  end

  def test_single_passing
    run = `./apgar-probe --document-root tmp --healthcheck-tree fixtures/001-single-passing`
    exitcode = $?.to_i
    assert_equal true, (exitcode == 0)
    assert_equal "OK\n", File.open(STATUS_FILE) { |file| file.read }
  end

end
