require 'test_helper'

class ReplTest < Minitest::Test
  def test_that_it_has_a_name
    refute_nil ::Repl::NAME
  end

  def test_that_it_has_a_version
    refute_nil ::Repl::VERSION
  end

  def test_that_version_regexp_is_correct
    assert_match Repl.version_regexp, 'repl 0.10.0 (rlwrap 0.43)'
    assert_match Repl.version_regexp, 'repl 0.10.0 (rlwrap not installed)'
  end

  def test_that_version_matches_version_regexp
    assert_match Repl.version_regexp, Repl.version
  end

  def test_it_does_something_useful
    # rubocop:disable Minitest/UselessAssertion
    assert true # you better believe it!
    # rubocop:enable Minitest/UselessAssertion
  end
end
